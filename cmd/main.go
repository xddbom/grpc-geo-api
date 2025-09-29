package main

import (
    "log"
    
    "grpc-geo-api/internal/app"
    "grpc-geo-api/internal/config"
)

func main() {
    cfg := config.Load()
    
    application, err := app.New(*cfg)
    if err != nil {
        log.Fatalf("Failed to initialize application: %v", err)
    }
    
 
    if err := application.Run(); err != nil {
        log.Fatalf("Application runtime error: %v", err)
    }
    
    log.Println("Application exited successfully")
}