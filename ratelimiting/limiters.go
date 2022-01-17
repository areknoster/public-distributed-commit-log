// Package ratelimiting provides rate limiting tools for sentinel.
package ratelimiting

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-limiter"
)

// TODO: this could be dynamic so it should be possible to do multi tenancy.
const storeKey = "key"

// TokenBucketLimiter does rate limiting based on token bucket algorithm.
type TokenBucketLimiter struct {
	store limiter.Store
}

func NewTokenBucketLimiter(store limiter.Store) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		store: store,
	}
}

func (l *TokenBucketLimiter) Limit() bool {
	tokens, rem, reset, ok, err := l.store.Take(context.Background(), storeKey)
	log.Debug().Msgf("all: %d, remaining: %d, reset: %v, ok: %v, err: %v", tokens, rem, reset, ok, err)
	if err != nil {
		log.Error().Err(err)
		return false
	}
	return !ok
}

// AlwaysAllowLimiter always passes requests.
type AlwaysAllowLimiter struct{}

func NewAlwaysAllowLimiter() *AlwaysAllowLimiter {
	return &AlwaysAllowLimiter{}
}

func (l *AlwaysAllowLimiter) Limit() bool {
	return false
}
