package token

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"icylight/uniswap/config"
	"icylight/uniswap/service"
)

// Token token result
type Token struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Logo     string `json:"logo"`
	Address  string `json:"address"`
}

// GetTokenByAddress get token by address
func GetTokenByAddress(address string) (*Token, error) {
	key := KeyToken()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	b, err := service.Redis.HGet(ctx, key, address).Bytes()
	if err == nil {
		var t Token
		json.Unmarshal(b, &t)
		return &t, nil
	}
	if err != redis.Nil {
		log.Println(err)
		return nil, err
	}

	raw, err := tokenMetaData(address)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(raw) == 0 {
		return nil, nil
	}

	var t Token
	err = json.Unmarshal(raw, &t)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	t.Address = address

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = service.Redis.HSet(ctx, key, address, string(raw)).Err()
	if err != nil {
		log.Println(err)
	}
	return &t, nil
}

// KeyToken redis 的 key (hash)
func KeyToken() string {
	return "token"
}

func tokenMetaData(address string) (json.RawMessage, error) {
	body := map[string]any{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "alchemy_getTokenMetadata",
		"params":  []string{address},
	}
	b, _ := json.Marshal(body)

	url := config.Conf.Eth.URL
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res struct {
		Result json.RawMessage `json:"result"`
	}
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}
