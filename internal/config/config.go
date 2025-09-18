package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"database"`

	Excel struct {
		Filepath string `yaml:"filepath"`
	} `yaml:"excel"`

	AI struct {
		APIKeyEnv   string  `yaml:"api_key_env"`
		BaseURL     string  `yaml:"base_url"`
		Model       string  `yaml:"model"`
		Temperature float64 `yaml:"temperature"`
		Enabled     bool    `yaml:"enabled"`
	}

	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
}

func LoadConfig(filepath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return config, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
	)
}

func (c *Config) GetAPIKey() string {
	return os.Getenv(c.AI.APIKeyEnv)
}
