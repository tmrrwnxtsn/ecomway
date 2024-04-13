package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const databaseURLEnvKey = "DATABASE_URL"

type Config struct {
	Engine   EngineConfig   `yaml:"engine"`
	Services ServicesConfig `yaml:"services"`
}

type EngineConfig struct {
	GRPCAddress string          `yaml:"grpc_address"`
	Storage     StorageConfig   `yaml:"storage"`
	Scheduler   SchedulerConfig `yaml:"scheduler"`
}

type StorageConfig struct {
	DatabaseURL string // подгружаем значения из переменных среды окружения
}

type SchedulerConfig struct {
	IsEnabled bool                 `yaml:"is_enabled"`
	Tasks     SchedulerTasksConfig `yaml:"tasks"`
}

type SchedulerTasksConfig struct {
	FinalizeOperations SchedulerTaskConfig `yaml:"finalize_operations"`
}

type SchedulerTaskConfig struct {
	IsEnabled                bool           `yaml:"is_enabled"`
	Interval                 int            `yaml:"interval"`
	OperationBatchSize       int64          `yaml:"operation_batch_size"`
	MaxWorkers               int            `yaml:"max_workers"`
	ActualizeStatusIntervals map[int]int    `yaml:"actualize_status_intervals"`
	ExternalSystemLifetime   map[string]int `yaml:"external_system_lifetime"`
}

type ServicesConfig struct {
	Integration ServiceConfig `yaml:"integration"`
}

type ServiceConfig struct {
	GRPCAddress string `yaml:"grpc_address"`
}

func Load(configPath string) (Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return cfg, err
	}

	cfg.loadFromEnv()

	return cfg, nil
}

func (c *Config) loadFromEnv() {
	if dsnFromEnv, exists := os.LookupEnv(databaseURLEnvKey); exists {
		c.Engine.Storage.DatabaseURL = dsnFromEnv
	}
}
