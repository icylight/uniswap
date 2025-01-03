package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Conf config
var Conf Config

func init() {
	Conf.Redis.Addr = "localhost:6379"
}

// Config global config
type Config struct {
	Redis Redis `yaml:"redis"`
	Eth   Eth   `yaml:"eth"`
}

// Redis config for redis
type Redis struct {
	Addr string `yaml:"addr"`
}

// Eth config
type Eth struct {
	URL string `yaml:"url"`
}

// Load config
func Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	d := yaml.NewDecoder(f)
	return d.Decode(&Conf)
}
