package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"icylight/uniswap/models/redis/pair"
	"icylight/uniswap/service"
)

type inputData struct {
	Params struct {
		Result struct {
			Transaction struct {
				Hash  string `json:"hash"`
				Input string `json:"input"`
			} `json:"transaction"`
		} `json:"result"`
	} `json:"params"`
}

// GetPairs 获取交易对
func GetPairs(message []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()
	var data inputData
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ok, err := service.Redis.SetNX(ctx, lockKeyPairs(data.Params.Result.Transaction.Hash), 1,
		5*time.Second).Result()
	if err != nil {
		log.Println(err)
		return
	}
	if !ok {
		return
	}
	if strings.HasPrefix(data.Params.Result.Transaction.Input, "0x38ed1739") {
		getPairsBySwapExactTokensForTokens(data.Params.Result.Transaction.Input)
		return
	}
	if strings.HasPrefix(data.Params.Result.Transaction.Input, "0x18cbafe5") {
		getPairBySwapExactTokensForETH(data.Params.Result.Transaction.Input)
		return
	}
	// TODO: 增加其他 swap 函数
}

// 通过 swapExactTokensForTokens 获取的 pairs
func getPairsBySwapExactTokensForTokens(input string) {
	getPairByMethodIndex(input, 2)
}

func getPairBySwapExactTokensForETH(input string) {
	getPairByMethodIndex(input, 2)
}

func getPairByMethodIndex(input string, addressIdx int) {
	hex := common.Hex2Bytes(input[2:])
	m, err := service.UniswapRouterABI.MethodById(hex[:4])
	if err != nil {
		log.Println(err)
		return
	}
	args, err := m.Inputs.Unpack(hex[4:])
	if err != nil {
		log.Println(err)
		return
	}
	if addressIdx > len(args) {
		log.Printf("addressIdx(%d) > len(args)(%d), method: %s\n", addressIdx, len(args), m.String())
		return
	}
	tokens, ok := args[addressIdx].([]common.Address)
	if !ok {
		log.Println("parse tokens failed")
		return
	}
	for i := 1; i < len(tokens); i++ {
		pair.HandlePair(pair.HandlerPairParam{
			Addr1: tokens[i-1],
			Addr2: tokens[i],
		})
	}
}

func lockKeyPairs(hash string) string {
	return fmt.Sprintf("lock:pairs:%s", hash)
}
