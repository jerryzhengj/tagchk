package util

import (
	cache "github.com/UncleBig/goCache"
	"time"
)

func NewCache(expir time.Duration, cleanup time.Duration) *cache.Cache{
	return cache.New(expir, cleanup)
}