package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	TenantRegistry TenantRegistry
	Sender         Sender
	Infrastructure Infrastructure
}

func Load() (*Config, error) {
	tenantRegistryCacheUpdateInterval, err := loadTenantRegistryCacheUpdateInterval()
	if err != nil {
		return nil, fmt.Errorf("failed to load receiver tenant cache interval: %w", err)
	}
	tenantRegistryEndpoint, err := loadTenantRegistryEndpoint()
	if err != nil {
		return nil, fmt.Errorf("failed to load tenant registry endpoint: %w", err)
	}
	sendEndpoint, err := loadSendEndpoint()
	if err != nil {
		return nil, fmt.Errorf("failed to load send endpoint: %w", err)
	}
	sendInterval, err := loadSendInterval()
	if err != nil {
		return nil, fmt.Errorf("failed to load send interval: %w", err)
	}
	sendConcurrency, err := loadSendConcurrency()
	if err != nil {
		return nil, fmt.Errorf("failed to load send concurrency: %w", err)
	}
	metricPort, err := loadMetricPort()
	if err != nil {
		return nil, fmt.Errorf("failed to load metric port: %w", err)
	}
	logLevel, err := loadLogLevel()
	if err != nil {
		return nil, fmt.Errorf("failed to load log level: %w", err)
	}
	return &Config{
		TenantRegistry: TenantRegistry{
			TenantRegistryEndpoint:            tenantRegistryEndpoint,
			TenantRegistryCacheUpdateInterval: tenantRegistryCacheUpdateInterval,
		},
		Sender: Sender{
			SendEndpoint:    sendEndpoint,
			SendInterval:    sendInterval,
			SendConcurrency: sendConcurrency,
		},
		Infrastructure: Infrastructure{
			MetricPort: metricPort,
			LogLevel:   logLevel,
		},
	}, nil
}

const (
	defaultTenantRegistryCacheUpdateInterval = time.Minute
)

type TenantRegistry struct {
	TenantRegistryEndpoint            string
	TenantRegistryCacheUpdateInterval time.Duration
}

func loadTenantRegistryCacheUpdateInterval() (time.Duration, error) {
	tenantRegistryCacheUpdateIntervalStr := os.Getenv("TENANT_REGISTRY_CACHE_UPDATE_INTERVAL")
	if tenantRegistryCacheUpdateIntervalStr == "" {
		return defaultTenantRegistryCacheUpdateInterval, nil
	}
	tenantRegistryCacheUpdateInterval, err := time.ParseDuration(tenantRegistryCacheUpdateIntervalStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse TENANT_REGISTRY_CACHE_UPDATE_INTERVAL: %w", err)
	}
	return tenantRegistryCacheUpdateInterval, nil
}

func loadTenantRegistryEndpoint() (string, error) {
	tenantRegistryEndpointStr := os.Getenv("TENANT_REGISTRY_ENDPOINT")
	if tenantRegistryEndpointStr == "" {
		return "", fmt.Errorf("environment variable TENANT_REGISTRY_ENDPOINT is not set")
	}
	return tenantRegistryEndpointStr, nil
}

const (
	defaultSendInterval    = 30 * time.Second
	defaultSendConcurrency = 10
)

type Sender struct {
	SendEndpoint    string
	SendInterval    time.Duration
	SendConcurrency int
}

func loadSendEndpoint() (string, error) {
	sendEndpointStr := os.Getenv("SEND_ENDPOINT")
	if sendEndpointStr == "" {
		return "", fmt.Errorf("SEND_ENDPOINT environment variable not set")
	}
	return sendEndpointStr, nil
}

func loadSendInterval() (time.Duration, error) {
	sendIntervalStr := os.Getenv("SEND_INTERVAL")
	if sendIntervalStr == "" {
		return defaultSendInterval, nil
	}
	sendInterval, err := time.ParseDuration(sendIntervalStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse SEND_INTERVAL: %w", err)
	}
	return sendInterval, nil
}

func loadSendConcurrency() (int, error) {
	sendConcurrencyStr := os.Getenv("SEND_CONCURRENCY")
	if sendConcurrencyStr == "" {
		return defaultSendConcurrency, nil
	}
	sendConcurrency, err := strconv.Atoi(sendConcurrencyStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse SEND_CONCURRENCY: %w", err)
	}
	return sendConcurrency, nil
}

const (
	defaultMetricPort = 10000
	defaultLogLevel   = slog.LevelInfo
)

type Infrastructure struct {
	MetricPort int
	LogLevel   slog.Level
}

func loadMetricPort() (int, error) {
	metricPortStr := os.Getenv("METRIC_PORT")
	if metricPortStr == "" {
		return defaultMetricPort, nil
	}
	metricPort, err := strconv.Atoi(metricPortStr)
	if err != nil {
		return 0, err
	}
	return metricPort, nil
}

func loadLogLevel() (slog.Level, error) {
	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		return defaultLogLevel, nil
	}
	var level slog.Level
	err := level.UnmarshalText([]byte(logLevelStr))
	if err != nil {
		return 0, err
	}
	return level, nil
}
