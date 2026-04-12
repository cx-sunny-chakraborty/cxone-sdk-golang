package auth

import (
	"context"
	"sync"
)

// fetcher is the unexported contract that concrete flows (OAuth, API key)
// implement. The cached wrapper takes care of serialization, generation
// tracking, and the lazy-first / refresh-on-stale semantics required by
// CLAUDE.md §4.5.
type fetcher interface {
	fetch(ctx context.Context) (string, error)
	close()
}

// cached wraps any fetcher with a serialized, generation-tracked token cache.
//
// All access to the cached token goes through a single sync.Mutex held for
// the duration of the underlying fetch. Holding the lock across the network
// call is intentional: it guarantees that concurrent callers serialize on
// one fetch and all observe the same fresh token, with no double-fetch on
// thundering-herd 401s.
type cached struct {
	fetcher fetcher

	mu    sync.Mutex
	token string
	gen   uint64
}

func newCached(f fetcher) *cached {
	return &cached{fetcher: f}
}

// Token implements [Authenticator].
func (c *cached) Token(ctx context.Context) (string, uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" {
		return c.token, c.gen, nil
	}
	return c.fetchLocked(ctx)
}

// Refresh implements [Authenticator]. If the cache has advanced past staleGen
// (i.e. another goroutine already refreshed) Refresh is a no-op.
func (c *cached) Refresh(ctx context.Context, staleGen uint64) (string, uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && c.gen > staleGen {
		return c.token, c.gen, nil
	}
	return c.fetchLocked(ctx)
}

// fetchLocked invokes the underlying fetcher with c.mu held. The lock is
// released by the caller.
func (c *cached) fetchLocked(ctx context.Context) (string, uint64, error) {
	tok, err := c.fetcher.fetch(ctx)
	if err != nil {
		return "", 0, err
	}
	c.token = tok
	c.gen++
	return c.token, c.gen, nil
}

// Close releases any resources held by the underlying fetcher and forgets the
// cached token. Subsequent calls to Token will return an error from the
// underlying fetcher rather than re-issuing IAM requests; the cxone client's
// Close path is the only expected caller.
func (c *cached) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = ""
	c.fetcher.close()
}
