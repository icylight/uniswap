package service

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"icylight/uniswap/config"
)

func TestMain(m *testing.M) {
	err := config.Load("../config.yml")
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func TestInit(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	Init()

	ctx := context.Background()
	require.NoError(Redis.Set(ctx, "test", "1", time.Second).Err())

	v, err := Redis.Get(ctx, "test").Result()
	require.NoError(err)
	assert.Equal("1", v)

	assert.NotNil(EthClient)
}
