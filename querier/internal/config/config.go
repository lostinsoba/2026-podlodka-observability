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
	Querier        Querier
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
	queryEndpoint, err := loadQueryEndpoint()
	if err != nil {
		return nil, fmt.Errorf("failed to load query endpoint: %w", err)
	}
	queryInterval, err := loadQueryInterval()
	if err != nil {
		return nil, fmt.Errorf("failed to load query interval: %w", err)
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
		Querier: Querier{
			QueryEndpoint: queryEndpoint,
			QueryInterval: queryInterval,
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
	defaultQueryInterval = time.Minute
)

type Querier struct {
	QueryEndpoint string
	QueryInterval time.Duration
}

func loadQueryEndpoint() (string, error) {
	queryEndpointStr := os.Getenv("QUERY_ENDPOINT")
	if queryEndpointStr == "" {
		return "", fmt.Errorf("QUERY_ENDPOINT environment variable not set")
	}
	return queryEndpointStr, nil
}

func loadQueryInterval() (time.Duration, error) {
	queryIntervalStr := os.Getenv("QUERY_INTERVAL")
	if queryIntervalStr == "" {
		return defaultQueryInterval, nil
	}
	queryInterval, err := time.ParseDuration(queryIntervalStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse QUERY_INTERVAL: %w", err)
	}
	return queryInterval, nil
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
