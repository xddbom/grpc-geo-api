package config 

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	GRPC		GRPCConfig
	Gateway		GatewayConfig
	Log			LogConfig
}

type GRPCConfig struct {
	Port	string 
	Host	string
}

type GatewayConfig struct {
	Port	string 
	Host	string
}

type LogConfig struct {
	Level		string
	Environment	string
}

func Load() *Config {
	return &Config{
		GRPC: GRPCConfig{
			Port: getEnv("GRPC_PORT", "50051"),
			Host: getEnv("GRPC_HOST", ""),
		},

		Gateway: GatewayConfig{
			Port: getEnv("GATEWAY_PORT", "8081"),
			Host: getEnv("GATEWAY_HOST", ""),
		},
		Log: LogConfig{
			Level:       getEnv("LOG_LEVEL", "debug"),
			Environment: getEnv("APP_ENV", "development"),
		},
	}
}


func (g GRPCConfig) Address() string {
	return g.Host + ":" + g.Port
}


func (g GatewayConfig) Address() string {
	return g.Host + ":" + g.Port
}


func getEnv[T any](key string, defaultValue T) T {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

    switch any(defaultValue).(type) {
    case string:
        return any(val).(T)
    case int:
        if v, err := strconv.Atoi(val); err == nil {
            return any(v).(T)
        }
    case time.Duration:
        if v, err := time.ParseDuration(val); err == nil {
            return any(v).(T)
        }
    }
    return defaultValue
}