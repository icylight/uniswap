package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	Init()

	ctx := context.Background()
	require.NoError(Redis.Set(ctx, "test", "1", time.Second).Err())

	v, err := Redis.Get(ctx, "test").Result()
	require.NoError(err)
	assert.Equal("1", v)
}
