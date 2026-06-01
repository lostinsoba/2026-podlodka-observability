package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	API            API
	TenantRegistry TenantRegistry
	Infrastructure Infrastructure
}

func Load() (*Config, error) {
	apiPort, err := loadAPIPort()
	if err != nil {
		return nil, fmt.Errorf("failed to load api port: %w", err)
	}
	tenantConfigs, err := loadTenantConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to load tenant config: %w", err)
	}
	metricPort, err := loadMetricPort()
	if err != nil {
		return nil, fmt.Errorf("failed to load telemetry port: %w", err)
	}
	logLevel, err := loadLogLevel()
	if err != nil {
		return nil, fmt.Errorf("failed to load log level: %w", err)
	}
	return &Config{
		API: API{
			Port: apiPort,
		},
		TenantRegistry: TenantRegistry{
			TenantConfigurations: tenantConfigs,
		},
		Infrastructure: Infrastructure{
			MetricPort: metricPort,
			LogLevel:   logLevel,
		},
	}, nil
}

const (
	defaultAPIPort = 8080
)

type API struct {
	Port int
}

func loadAPIPort() (int, error) {
	apiPortStr := os.Getenv("API_PORT")
	if apiPortStr == "" {
		return defaultAPIPort, nil
	}
	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil {
		return 0, err
	}
	return apiPort, nil
}

const (
	tenantConfigPath = "/etc/tenant-registry/tenants.yml"
)

type TenantRegistry struct {
	TenantConfigurations TenantConfigurations
}

type TenantConfigurations struct {
	Tenants []TenantConfiguration `yaml:"tenants"`
}

type TenantConfiguration struct {
	TenantID           string `yaml:"tenant_id"`
	TenantMaxBatchSize int    `yaml:"tenant_max_batch_size"`
}

func loadTenantConfigs() (TenantConfigurations, error) {
	var cfg TenantConfigurations
	data, err := os.ReadFile(tenantConfigPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read tenant config file: %w", err)
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to unmarshal tenant config: %w", err)
	}
	return cfg, nil
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
