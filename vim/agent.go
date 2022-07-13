package main

import (
	"encoding/gob"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/zgs225/youdao"
)

const (
	CACHE_EXPIRES time.Duration = 30 * 24 * time.Hour
)

type agentClient struct {
	Client *youdao.Client
	Cache  *cache.Cache
	Dirty  bool
}

func (a *agentClient) Query(q string) (*youdao.Result, error) {
	r, err := a.Client.Query(q)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func newAgent(c *youdao.Client) *agentClient {
	gob.Register(&youdao.Result{})
	return &agentClient{c, nil, false}
}
