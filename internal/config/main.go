package config

import (
	"os"

	pbTypes "github.com/m1dugh/program-browser/pkg/types"
	"gopkg.in/yaml.v3"
)

type redisConfig struct {
	Addr              string `yaml:"address"`
	PlainTextPassword string `yaml:"password"`
	DB                int    `yaml:"db"`
	Name              string `yaml:"name"`
}

type fileConfig struct {
	Format string `yaml:"format"`
	Target string `yaml:"filename,omitempty"`
}

type OutputConfig struct {
	Redis *redisConfig `yaml:"redis"`
	File  *fileConfig  `yaml:"file"`
}

type bugcrowdConfig struct {
	Enable bool `yaml:"enable"`
}

type InputConfig struct {
	/// A wildcard filter based on program names
	Filters []NameFilter `yaml:"filters"`

	/// The Bugcrowd config
	Bugcrowd *bugcrowdConfig `yaml:"bugcrowd"`

	/// Extra entries of program to push to output
	ExtraEntries []pbTypes.Program `yaml:"extraEntries"`
}

type Config struct {
	Output OutputConfig `yaml:"output"`
	Input  InputConfig  `yaml:"input"`
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

	for _, filter := range config.Input.Filters {
		filter.configure()
	}

	return config, nil
}
