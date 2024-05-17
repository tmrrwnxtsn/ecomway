package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const apiKeyEnvKey = "API_KEY"

type Config struct {
	Report   ReportConfig   `yaml:"report"`
	Services ServicesConfig `yaml:"services"`
}

type ReportConfig struct {
	HTTPAddress string `yaml:"http_address"`
	APIKey      string
}

type ServicesConfig struct {
	Engine ServiceConfig `yaml:"engine"`
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
	if apiKey, exists := os.LookupEnv(apiKeyEnvKey); exists {
		c.Report.APIKey = apiKey
	}
}
