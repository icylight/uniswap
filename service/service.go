package service

import (
	"icylight/uniswap/config"

	"github.com/redis/go-redis/v9"
)

// services
var (
	Redis *redis.Client
)

// Init init
func Init() error {
	r := config.Conf.Redis

	Redis = redis.NewClient(&redis.Options{
		Addr: r.Addr,
	})
	return nil
}
