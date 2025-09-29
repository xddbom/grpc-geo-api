package app

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "sync"

    "grpc-geo-api/internal/config"
)

type App struct {
    grpcServer *Server
    gateway    *Gateway
    config     config.Config
}

func New(cfg config.Config) (*App, error) {
    grpcServer, err := NewServer(cfg.GRPC)
    if err != nil {
        return nil, fmt.Errorf("failed to create gRPC server: %w", err)
    }

    grpcServer.RegisterServices()

    gateway, err := NewGateway(cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to create gateway: %w", err)
    }

    return &App{
        grpcServer: grpcServer,
        gateway:    gateway,
        config:     cfg,
    }, nil
}

func (a *App) Run() error {
    errChan := make(chan error, 2)
    
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        log.Printf("Starting gRPC server on %s", a.config.GRPC.Address())
        if err := a.grpcServer.Start(); err != nil {
            errChan <- fmt.Errorf("gRPC server error: %w", err)
        }
    }()

    go func() {
        defer wg.Done()
        log.Printf("Starting HTTP gateway on %s", a.config.Gateway.Address())
        if err := a.gateway.Start(); err != nil {
            errChan <- fmt.Errorf("gateway error: %w", err)
        }
    }()

    return a.waitForShutdown(errChan, &wg)
}

func (a *App) waitForShutdown(errChan chan error, wg *sync.WaitGroup) error {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    select {
    case err := <-errChan:
        log.Printf("Server error occurred: %v", err)
        a.Stop()
        return err
    case sig := <-sigChan:
        log.Printf("Received signal: %v", sig)
        a.Stop()
        wg.Wait() 
        return nil
    }
}

func (a *App) Stop() {
    log.Println("Shutting down application...")

    log.Println("Stopping gateway...")
    a.gateway.Stop()

    log.Println("Stopping gRPC server...")
    a.grpcServer.Stop()

    log.Println("Application stopped gracefully")
}