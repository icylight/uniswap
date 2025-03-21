package token

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"icylight/uniswap/config"
	"icylight/uniswap/service"
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

func TestGetTokenByAddress(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	addr := "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"

	token, err := GetTokenByAddress(addr)
	require.NoError(err)
	require.NotNil(token)
	assert.Equal("WETH", token.Name)

	// 缓存
	token, err = GetTokenByAddress(addr)
	require.NoError(err)
	require.NotNil(token)
	assert.Equal("WETH", token.Name)

	cache, err := service.Redis.HGet(context.TODO(), KeyToken(), addr).Result()
	require.NoError(err)
	t.Log(cache)
}
