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
	Redis      Redis      `yaml:"redis"`
	Eth        Eth        `yaml:"eth"`
	UniSwapABI UniSwapABI `yaml:"uniswap_abi"`
}

// Redis config for redis
type Redis struct {
	Addr string `yaml:"addr"`
}

// Eth config
type Eth struct {
	URL string `yaml:"url"`
	WS  string `yaml:"ws"`
}

// UniSwapABI uniswap abi json files
type UniSwapABI struct {
	PairABI   string `yaml:"pair_abi"`
	RouterABi string `yaml:"router_abi"`
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
