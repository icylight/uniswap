package pair

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strconv"
	"time"

	"icylight/uniswap/models/redis/token"
	"icylight/uniswap/service"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/miraclesu/uniswap-sdk-go/constants"
	"github.com/redis/go-redis/v9"
)

// FormatPair 将交易对排序
func FormatPair(tokens []string) {
	sort.Strings(tokens)
}

// KeyPair 交易对缓存 key zset
func KeyPair(tokens []string) string {
	return fmt.Sprintf("pair:%s:%s", tokens[0], tokens[1])
}

// KeyLatestPair 最新的交易对缓存 hash
func KeyLatestPair() string {
	return "latest_pair"
}

// HandlerPairParam 处理交易对的参数
type HandlerPairParam struct {
	Addr1 common.Address
	Addr2 common.Address
}

// HandlePair 处理交易对
// TODO: 参数需要带上时间戳，避免旧数据覆盖新数据
func HandlePair(param HandlerPairParam) {
	if param.Addr1.String() > param.Addr2.String() {
		param.Addr1, param.Addr2 = param.Addr2, param.Addr1
	}
	token0, err := token.GetTokenByAddress(param.Addr1.String())
	if err != nil {
		log.Println(err)
		return
	}
	if token0 == nil {
		log.Println("token0 not exist")
		return
	}
	token1, err := token.GetTokenByAddress(param.Addr2.String())
	if err != nil {
		log.Println(err)
		return
	}
	if token1 == nil {
		log.Println("token1 not exist")
		return
	}

	contractAddr := getCreate2Address(param.Addr1, param.Addr2)

	// REVIEW: calldata 是否可以放入本地缓存
	callData, err := service.UniswapPairABI.Pack("getReserves")
	if err != nil {
		log.Println(err)
		return
	}
	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: callData,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := service.EthClient.CallContract(ctx, msg, nil)
	if err != nil {
		log.Println(err)
		return
	}

	var reserve0, reserve1 *big.Int
	var timestamp uint32
	err = service.UniswapPairABI.UnpackIntoInterface(&[]any{&reserve0, &reserve1, &timestamp},
		"getReserves", result)
	if err != nil {
		log.Println(err)
		return
	}
	price := new(big.Float).Quo(new(big.Float).SetInt(reserve1), new(big.Float).SetInt(reserve0))

	p := Pair{
		Token0:    token0,
		Token1:    token1,
		Reserve0:  reserve0.String(),
		Reserve1:  reserve1.String(),
		Price:     price.String(),
		Timestamp: strconv.FormatInt(int64(timestamp), 10),
	}

	bytes, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
		return
	}

	key := KeyPair([]string{p.Token0.Symbol, p.Token1.Symbol})
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = service.Redis.ZAdd(ctx, key, redis.Z{
		Score:  float64(timestamp),
		Member: bytes,
	}).Err()
	if err != nil {
		log.Println(err)
		return
	}

	s, err := service.Redis.ZRevRange(ctx, key, 0, 0).Result()
	if err != nil {
		log.Println(err)
		return
	}

	if len(s) > 0 {
		err = service.Redis.HSet(ctx, KeyLatestPair(), key, s[0]).Err()
		if err != nil {
			log.Println(err)
			return
		}
	}
	fmt.Printf("update %s\n", key)
}

func getCreate2Address(addressA, addressB common.Address) common.Address {
	var salt [32]byte
	copy(salt[:], crypto.Keccak256(append(addressA.Bytes(), addressB.Bytes()...)))
	return crypto.CreateAddress2(constants.FactoryAddress, salt, constants.InitCodeHash)
}

// Pair redis pair
type Pair struct {
	Token0    *token.Token `json:"token0"`
	Token1    *token.Token `json:"token1"`
	Reserve0  string       `json:"reserve0"`
	Reserve1  string       `json:"reserve1"`
	Price     string       `json:"price"`
	Timestamp string       `json:"timestamp"`
}
