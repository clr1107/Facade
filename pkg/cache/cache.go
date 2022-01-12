package cache

import (
	gocache "github.com/patrickmn/go-cache"
	"time"
)

type cacheUnit struct {
	data    interface{}
	isValid func(u *cacheUnit) bool
}

// NewCacheUnit creates a new cacheable chunk of data. Optionally a function
// `isValid` can be provided to check if a `*cacheUnit` is valid. If this
// function fails the cache will purge the key.
func NewCacheUnit(data interface{}, isValid func(u *cacheUnit) bool) *cacheUnit {
	return &cacheUnit{
		data:    data,
		isValid: isValid,
	}
}

// Cacher represents all cache implementations, barebones interoperable
// functions, for any implementation that may be used down the line.
type Cacher interface {
	// Put sets, if not present, the key to the `*cache.Unit`, with default
	// expiry. May return an error depending on impl.
	Put(key string, unit *cacheUnit) error
	// PutExpiry sets, if not present, the key to the `*cache.Unit`, with
	// given expiry. May return an error depending on impl.
	PutExpiry(key string, unit *cacheUnit, expiry time.Time) error
	// Set sets, the key to the `*cache.Unit`, present or not, with default
	// expiry. May return an error depending on impl.
	Set(key string, unit *cacheUnit) error
	// SetExpiry sets, the key to the `*cache.Unit`, present or not, with
	// given expiry. May return an error depending on impl.
	SetExpiry(key string, unit *cacheUnit, expiry time.Time) error
    // Get returns the key's unit's value if valid. Checks the `isValid`
	// function & expiry times. Bool indicates if present (& valid).
	Get(key string) (interface{}, bool)
	// GetDefault returns the key's unit's value if valid. Checks the `isValid`
	// function & expiry times. If not present (or valid) `d` is returned.
	GetDefault(key string, d interface{}) interface{}
    // Delete will remove the key from the cache. May return an error depending
	// on impl.
	Delete(key string) error
	Clear()
}

// CacheImpl implements patrickmn's `go-cache`.
type CacheImpl struct {
	cache *gocache.Cache
}

func NewCache(defaultTtl *time.Duration, defaultClean *time.Duration) *CacheImpl {
	var ttl, clean time.Duration

	if defaultTtl == nil {
		ttl = gocache.NoExpiration
	}

	if defaultClean == nil {
		clean = time.Minute * 15
	}

	return &CacheImpl{
		cache: gocache.New(ttl, clean),
	}
}

func (c *CacheImpl) Put(key string, unit *cacheUnit) error {
	return c.cache.Add(key, unit, gocache.DefaultExpiration)
}

func (c *CacheImpl) PutExpiry(key string, unit *cacheUnit, expiry time.Time) error {
	return c.cache.Add(key, unit, expiry.Sub(time.Now()))
}

func (c *CacheImpl) Set(key string, unit *cacheUnit) error {
	c.cache.Set(key, unit, gocache.DefaultExpiration)
	return nil
}

func (c *CacheImpl) SetExpiry(key string, unit *cacheUnit, expiry time.Time) error {
	c.cache.Set(key, unit, expiry.Sub(time.Now()))
	return nil
}

func (c *CacheImpl) Get(key string) (interface{}, bool) {
	var ret interface{}

	unit, found := c.cache.Get(key)
	if !found || unit == nil {
		return nil, false
	}

	cast, ok := unit.(*cacheUnit)
	if ok {
		if cast.isValid != nil && !cast.isValid(cast) {
			_ = c.Delete(key)
			return nil, false
		}

		ret = cast.data
	}

	return ret, ok
}

func (c *CacheImpl) GetDefault(key string, d interface{}) interface{} {
	ret, ok := c.Get(key)
	if !ok {
		ret = d
	}

	return ret
}

func (c *CacheImpl) Delete(key string) error {
	c.cache.Delete(key)
	return nil
}

func (c *CacheImpl) Clear() {
	c.cache.Flush()
}
