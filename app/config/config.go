package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database  DatabaseConfig  `yaml:"database"`
	API       APIConfig       `yaml:"api"`
	Simulator SimulatorConfig `yaml:"simulator"`
}

type DatabaseConfig struct {
	PrimaryDSN string `yaml:"primary_dsn"`
	ReplicaDSN string `yaml:"replica_dsn"`
}

type APIConfig struct {
	Port int `yaml:"port"`
}

type SimulatorConfig struct {
	Enabled       bool           `yaml:"enabled"`
	RatePerSecond int            `yaml:"rate_per_second"`
	Workers       int            `yaml:"workers"`
	Weights       WeightConfig   `yaml:"weights"`
}

type WeightConfig struct {
	Rental         int `yaml:"rental"`
	Return         int `yaml:"return"`
	CustomerChurn  int `yaml:"customer_churn"`
	Read           int `yaml:"read"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
