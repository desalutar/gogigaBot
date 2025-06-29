package config

import "os"

type Config struct {
    Env       string       `env:"ENV" envDefault:"development"`
    Server    ServerConfig
    Logger    LoggerConfig
    Cache     CacheConfig
}

type ServerConfig struct {
    Port string `env:"SERVER_PORT"`
}

type LoggerConfig struct {
    Level       string `env:"LOG_LEVEL" envDefault:"info"`
    Environment string `env:"ENV" envDefault:"development"`
}

type CacheConfig struct {
    Host string    `env:"CACHE_HOST" envDefault:"localhost"`
    Port string    `env:"CACHE_PORT" envDefault:"6379"`
}

func getenvOrDefault(key, def string) string {
    val := os.Getenv(key)
    if val == "" {
        return def
    }
    return val
}

func LoadConfig() *Config {
    return &Config{
        Env: getenvOrDefault("ENV", "development"),
        Server: ServerConfig{
            Port: getenvOrDefault("SERVER_PORT", ""),
        },
        Logger: LoggerConfig{
            Level:       getenvOrDefault("LOG_LEVEL", "info"),
            Environment: getenvOrDefault("LOG_ENVIRONMENT", "development"),
        },
        Cache: CacheConfig{
            Host: getenvOrDefault("CACHE_HOST", "localhost"),
            Port: getenvOrDefault("CACHE_PORT", "6379"),
        },
    }
}
