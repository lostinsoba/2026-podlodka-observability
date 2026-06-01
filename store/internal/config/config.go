package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	API            API
	TenantRegistry TenantRegistry
	Receiver       Receiver
	Database       Database
	Infrastructure Infrastructure
}

func Load() (*Config, error) {
	databaseConnStr, err := loadDatabaseConnStr()
	if err != nil {
		return nil, fmt.Errorf("failed to load database connection string: %w", err)
	}
	databaseMaxOpenConns, err := loadDatabaseMaxOpenConns()
	if err != nil {
		return nil, fmt.Errorf("failed to load database max open connections: %w", err)
	}
	databaseConnMaxLifetime, err := loadConnMaxLifetime()
	if err != nil {
		return nil, fmt.Errorf("failed to load conn max lifetime: %w", err)
	}
	receiverBatchSize, err := loadReceiverBatchSize()
	if err != nil {
		return nil, fmt.Errorf("failed to load receiver batch size: %w", err)
	}
	receiverQueueCapacity, err := loadReceiverQueueCapacity()
	if err != nil {
		return nil, fmt.Errorf("failed to load receiver queue capacity: %w", err)
	}
	receiverBatchProcessInterval, err := loadReceiverBatchProcessInterval()
	if err != nil {
		return nil, fmt.Errorf("failed to load receiver batch process interval: %w", err)
	}
	tenantRegistryCacheUpdateInterval, err := loadTenantRegistryCacheUpdateInterval()
	if err != nil {
		return nil, fmt.Errorf("failed to load receiver tenant cache interval: %w", err)
	}
	tenantRegistryEndpoint, err := loadTenantRegistryEndpoint()
	if err != nil {
		return nil, fmt.Errorf("failed to load tenant registry endpoint: %w", err)
	}
	apiPort, err := loadAPIPort()
	if err != nil {
		return nil, fmt.Errorf("failed to load api port: %w", err)
	}
	metricPort, err := loadMetricPort()
	if err != nil {
		return nil, fmt.Errorf("failed to load telemetry port: %w", err)
	}
	traceExportEndpoint := loadTraceExportEndpoint()
	traceSpanMaxQueueSize, err := loadTraceSpanMaxQueueSize()
	if err != nil {
		return nil, fmt.Errorf("failed to load trace span max queue size: %w", err)
	}
	traceSpanMaxExportBatchSize, err := loadTraceSpanMaxExportBatchSize()
	if err != nil {
		return nil, fmt.Errorf("failed to load trace span max export batch size: %w", err)
	}
	traceSpanBatchTimeout, err := loadTraceSpanBatchTimeout()
	if err != nil {
		return nil, fmt.Errorf("failed to load trace span batch timeout: %w", err)
	}
	traceSampleRate, err := loadTraceSampleRate()
	if err != nil {
		return nil, fmt.Errorf("failed to load trace sample rate: %w", err)
	}
	logLevel, err := loadLogLevel()
	if err != nil {
		return nil, fmt.Errorf("failed to load log level: %w", err)
	}
	return &Config{
		API: API{
			Port: apiPort,
		},
		Database: Database{
			ConnStr:         databaseConnStr,
			MaxOpenConns:    databaseMaxOpenConns,
			ConnMaxLifetime: databaseConnMaxLifetime,
		},
		Receiver: Receiver{
			ReceiverBatchSize:            receiverBatchSize,
			ReceiverQueueCapacity:        receiverQueueCapacity,
			ReceiverBatchProcessInterval: receiverBatchProcessInterval,
		},
		TenantRegistry: TenantRegistry{
			TenantRegistryEndpoint:            tenantRegistryEndpoint,
			TenantRegistryCacheUpdateInterval: tenantRegistryCacheUpdateInterval,
		},
		Infrastructure: Infrastructure{
			MetricPort:                  metricPort,
			TraceExportEndpoint:         traceExportEndpoint,
			TraceSpanMaxQueueSize:       traceSpanMaxQueueSize,
			TraceSpanMaxExportBatchSize: traceSpanMaxExportBatchSize,
			TraceSpanBatchTimeout:       traceSpanBatchTimeout,
			TraceSampleRate:             traceSampleRate,
			LogLevel:                    logLevel,
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
	defaultReceiverBatchSize            = 100
	defaultReceiverQueueCapacity        = 10000
	defaultReceiverBatchProcessInterval = time.Minute
)

type Receiver struct {
	ReceiverBatchSize            int
	ReceiverQueueCapacity        int
	ReceiverBatchProcessInterval time.Duration
}

func loadReceiverBatchSize() (int, error) {
	receiverBatchSizeStr := os.Getenv("RECEIVER_BATCH_SIZE")
	if receiverBatchSizeStr == "" {
		return defaultReceiverBatchSize, nil
	}
	receiverBatchSize, err := strconv.Atoi(receiverBatchSizeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse RECEIVER_BATCH_SIZE: %w", err)
	}
	return receiverBatchSize, nil
}

func loadReceiverQueueCapacity() (int, error) {
	receiverQueueCapacityStr := os.Getenv("RECEIVER_QUEUE_CAPACITY")
	if receiverQueueCapacityStr == "" {
		return defaultReceiverQueueCapacity, nil
	}
	receiverQueueCapacity, err := strconv.Atoi(receiverQueueCapacityStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse RECEIVER_QUEUE_CAPACITY: %w", err)
	}
	return receiverQueueCapacity, nil
}

func loadReceiverBatchProcessInterval() (time.Duration, error) {
	receiverBatchProcessIntervalStr := os.Getenv("RECEIVER_BATCH_PROCESS_INTERVAL")
	if receiverBatchProcessIntervalStr == "" {
		return defaultReceiverBatchProcessInterval, nil
	}
	receiverBatchProcessInterval, err := time.ParseDuration(receiverBatchProcessIntervalStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse RECEIVER_BATCH_PROCESS_INTERVAL: %w", err)
	}
	return receiverBatchProcessInterval, nil
}

const (
	defaultDatabaseMaxOpenConns    = 5
	defaultDatabaseConnMaxLifetime = time.Minute
)

type Database struct {
	ConnStr         string
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func loadDatabaseConnStr() (string, error) {
	databaseConnStr := os.Getenv("DB_CONN_STR")
	if databaseConnStr == "" {
		return "", fmt.Errorf("environment variable DB_CONN_STR is not set")
	}
	return databaseConnStr, nil
}

func loadDatabaseMaxOpenConns() (int, error) {
	databaseMaxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS")
	if databaseMaxOpenConnsStr == "" {
		return defaultDatabaseMaxOpenConns, nil
	}
	databaseMaxOpenConns, err := strconv.Atoi(databaseMaxOpenConnsStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse DB_MAX_OPEN_CONNS: %w", err)
	}
	return databaseMaxOpenConns, nil
}

func loadConnMaxLifetime() (time.Duration, error) {
	connMaxLifetimeStr := os.Getenv("CONN_MAX_LIFETIME")
	if connMaxLifetimeStr == "" {
		return defaultDatabaseConnMaxLifetime, nil
	}
	connMaxLifetime, err := time.ParseDuration(connMaxLifetimeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse CONN_MAX_LIFETIME: %w", err)
	}
	return connMaxLifetime, nil
}

const (
	defaultMetricPort = 10000
	defaultLogLevel   = slog.LevelInfo
)

type Infrastructure struct {
	MetricPort                  int
	TraceExportEndpoint         string
	TraceSpanMaxQueueSize       int
	TraceSpanMaxExportBatchSize int
	TraceSpanBatchTimeout       time.Duration
	TraceSampleRate             float64
	LogLevel                    slog.Level
}

func loadMetricPort() (int, error) {
	metricPortStr := os.Getenv("METRIC_PORT")
	if metricPortStr == "" {
		return defaultMetricPort, nil
	}
	metricPort, err := strconv.Atoi(metricPortStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse METRIC_PORT: %w", err)
	}
	return metricPort, nil
}

const (
	defaultTraceSpanMaxQueueSize       = 2048
	defaultTraceSpanMaxExportBatchSize = 512
	defaultTraceSpanBatchTimeout       = 5 * time.Second
	defaultTraceSampleRate             = 1
)

func loadTraceExportEndpoint() string {
	return os.Getenv("TRACE_EXPORT_ENDPOINT")
}

func loadTraceSpanMaxQueueSize() (int, error) {
	traceSpanMaxQueueSizeStr := os.Getenv("TRACE_SPAN_MAX_QUEUE_SIZE")
	if traceSpanMaxQueueSizeStr == "" {
		return defaultTraceSpanMaxQueueSize, nil
	}
	traceSpanMaxQueueSize, err := strconv.Atoi(traceSpanMaxQueueSizeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse TRACE_SPAN_MAX_QUEUE_SIZE: %w", err)
	}
	return traceSpanMaxQueueSize, nil
}

func loadTraceSpanMaxExportBatchSize() (int, error) {
	traceSpanMaxExportBatchSizeStr := os.Getenv("TRACE_SPAN_MAX_EXPORT_BATCH_SIZE")
	if traceSpanMaxExportBatchSizeStr == "" {
		return defaultTraceSpanMaxExportBatchSize, nil
	}
	traceSpanMaxExportBatchSize, err := strconv.Atoi(traceSpanMaxExportBatchSizeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse TRACE_SPAN_MAX_EXPORT_BATCH_SIZE: %w", err)
	}
	return traceSpanMaxExportBatchSize, nil
}

func loadTraceSpanBatchTimeout() (time.Duration, error) {
	traceSpanBatchTimeoutStr := os.Getenv("TRACE_SPAN_BATCH_TIMEOUT")
	if traceSpanBatchTimeoutStr == "" {
		return defaultTraceSpanBatchTimeout, nil
	}
	traceSpanBatchTimeout, err := time.ParseDuration(traceSpanBatchTimeoutStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse TRACE_SPAN_BATCH_TIMEOUT: %w", err)
	}
	return traceSpanBatchTimeout, nil
}

func loadTraceSampleRate() (float64, error) {
	traceSampleRateStr := os.Getenv("TRACE_SAMPLE_RATE")
	if traceSampleRateStr == "" {
		return defaultTraceSampleRate, nil
	}
	traceSampleRate, err := strconv.ParseFloat(traceSampleRateStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse TRACE_SAMPLE_RATE: %w", err)
	}
	return traceSampleRate, nil
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
