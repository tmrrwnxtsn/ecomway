package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Engine   EngineConfig   `yaml:"engine"`
	Services ServicesConfig `yaml:"services"`
}

type EngineConfig struct {
	GRPCAddress string `yaml:"grpc_address"`
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

	return cfg, nil
}
