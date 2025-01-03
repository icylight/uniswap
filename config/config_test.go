package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(Load("../config.example.yml"))
	assert.True(Conf.Eth.URL != "")
	t.Log(Conf.Eth.URL)
}
