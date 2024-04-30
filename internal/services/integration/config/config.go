package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Integration IntegrationConfig `yaml:"integration"`
	Services    ServicesConfig    `yaml:"services"`
}

type IntegrationConfig struct {
	GRPCAddress string          `yaml:"grpc_address"`
	YooKassa    *YooKassaConfig `yaml:"yookassa"`
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
	if key, exists := os.LookupEnv(yooKassaPaymentsSecretKeyEnvKey); exists {
		c.Integration.YooKassa.API.Payments.SecretKey = key
	}
	if key, exists := os.LookupEnv(yooKassaPayoutsSecretKeyEnvKey); exists {
		c.Integration.YooKassa.API.Payouts.SecretKey = key
	}
}
