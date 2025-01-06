package pair

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"icylight/uniswap/config"
	"icylight/uniswap/service"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	err := config.Load("../../../config.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = service.Init()
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func TestGetReverse(t *testing.T) {
	require := require.New(t)

	addr1 := "0xB370b0E268d011d2758e8A8dE30028e7e74B2D24"

	addrWeth := "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	addrNewTech := "0xd0e6d04c2f105344860d07912a857ad21204fc97"
	fmt.Println(addrWeth < addrNewTech)

	address := common.HexToAddress(addr1)

	f, err := os.Open("/home/wangsen/go/src/icylight/uniswap/data/uniswap_v2_pair.abi.json")
	require.NoError(err)

	pairABI, err := abi.JSON(f)
	require.NoError(err)

	callData, err := pairABI.Pack("getReserves")
	require.NoError(err)
	msg := ethereum.CallMsg{
		To:   &address,
		Data: callData,
	}
	result, err := service.EthClient.CallContract(context.Background(), msg, nil)
	require.NoError(err)

	v, err := pairABI.Unpack("getReserves", result)
	require.NoError(err)
	fmt.Println(v)
	for i := range v {
		fmt.Println(reflect.TypeOf(v[i]))
	}
}

func TestHandlePair(t *testing.T) {
	require := require.New(t)
	param := HandlerPairParam{
		Addr1: common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
		Addr2: common.HexToAddress("0xd0e6d04c2f105344860d07912a857ad21204fc97"),
	}

	HandlePair(param)

	ctx := context.Background()

	v, err := service.Redis.HGetAll(ctx, KeyLatestPair()).Result()
	require.NoError(err)
	t.Log(v)
}
