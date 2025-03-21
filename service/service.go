package service

import (
	"log"
	"os"

	"icylight/uniswap/config"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
)

// services
var (
	Redis *redis.Client

	EthClient *ethclient.Client

	UniswapRouterABI abi.ABI
	UniswapPairABI   abi.ABI
)

// Init init
func Init() error {
	r := config.Conf.Redis

	Redis = redis.NewClient(&redis.Options{
		Addr: r.Addr,
	})

	var err error
	EthClient, err = ethclient.Dial(config.Conf.Eth.URL)
	if err != nil {
		log.Println(err)
		return err
	}

	f, err := os.Open(config.Conf.UniSwapABI.RouterABi)
	if err != nil {
		log.Println(err)
		return err
	}
	UniswapRouterABI, err = abi.JSON(f)
	if err != nil {
		log.Println(err)
		return err
	}

	f, err = os.Open(config.Conf.UniSwapABI.PairABI)
	if err != nil {
		log.Println(err)
		return err
	}
	UniswapPairABI, err = abi.JSON(f)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
