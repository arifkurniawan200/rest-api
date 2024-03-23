package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Database struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Name       string `yaml:"name"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	ActivePool bool   `yaml:"active_pool"`
	MaxPool    int    `yaml:"max_pool"`
	MinPool    int    `yaml:"min_pool"`
}

type Config struct {
	DB Database `yaml:"db"`
}

func ReadConfig() Config {
	data, err := os.ReadFile("config/app.yaml")
	if err != nil {
		panic(err)
	}

	var configuration Config

	if err := yaml.Unmarshal(data, &configuration); err != nil {
		panic(err)
	}
	return configuration
}
