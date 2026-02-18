package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type redisConfig struct {
	Addr string `yaml:"address"`
	PlainTextPassword string `yaml:"password"`
	DB int `yaml:"db"`
	Name string `yaml:"name"`
}

type Config struct {
	Redis *redisConfig `yaml:"redis"`
}

func (r redisConfig) Password() string {
	if r.PlainTextPassword != "" {
		return r.PlainTextPassword
	}

	return os.Getenv("REDIS_PASSWORD")
}

func NewConfig(configFile string) (Config, error) {

	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
