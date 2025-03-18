package cache

import (
	"context"
	"fmt"
)

type CacheType int

const (
	Redis     CacheType = iota
	Dice      CacheType = iota
	redisPort           = 6379
	dicePort            = 7379
	cacheHost           = "localhost"
)

type Cache interface {
	Stream(ctx context.Context, key string, ch chan<- any)
}

func (c CacheType) String() string {
	return [...]string{"redis", "dice"}[c]
}

func NewCache(cache CacheType) (Cache, error) {
	switch cache {
	case Redis:
		return newRedisCache(fmt.Sprintf("%s:%d", cacheHost, redisPort))
	case Dice:
		return newDiceCache(cacheHost, dicePort)
	}
	return nil, fmt.Errorf("unknown cache type (%v) requested", cache)
}
