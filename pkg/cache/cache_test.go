package cache

import (
	"testing"
	"time"
)

// Some implementation functions don't actually return errors, despite indicating
// Due to how patricknm's library is written. So that's why some are ignored.

func TestCacheImpl_Delete(t *testing.T) {
	c := NewCache(nil, nil)
	_ = c.Set("key", NewCacheUnit(5, nil))

	_, present := c.Get("key")
	if !present {
		t.Errorf("cache unit not present after adding")
	}

	_ = c.Delete("key")
	if _, present2 := c.Get("key"); present2 {
		t.Errorf("cache unit present after deleting")
	}
}

func TestCacheImpl_Get(t *testing.T) {
	c := NewCache(nil, nil)

	_ = c.Set("one", NewCacheUnit(1, nil))
	if one, present := c.Get("one"); !present {
		t.Errorf("cache unit not preent after adding")
	} else if one != 1 {
		t.Errorf("cache unit value not correct, wanted %v got %v", 1, one)
	}

	if _, present := c.Get("two"); present {
		t.Errorf("cache unit apparently present when it does not exist")
	}
}

func TestCacheImpl_GetDefault(t *testing.T) {
	c := NewCache(nil, nil)

	if val := c.GetDefault("key", 5); val != 5 {
		t.Errorf("value from GetDefault was %v wanted %v", val, 5)
	}
}

func TestCacheImpl_Put(t *testing.T) {
	c := NewCache(nil, nil)
	_ = c.Put("key", NewCacheUnit(2, nil)) // impl never errors

	if val, ok := c.Get("key"); !ok {
		t.Errorf("cache unit not present after adding")
	} else if val != 2 {
		t.Errorf("value after Put is incorrect")
	}
}

func TestCacheImpl_PutExpiry(t *testing.T) {
	c := NewCache(nil, nil)
	_ = c.PutExpiry("key", NewCacheUnit(2, nil), time.Now().Add(time.Millisecond*250)) // impl never errors

	if val, ok := c.Get("key"); !ok || val != 2 {
		t.Errorf("cache unit not (correctly) present immediately after adding, wanted %v got %v", 2, val)
	}

	time.Sleep(time.Millisecond * 300)
	if _, ok := c.Get("key"); ok {
		t.Errorf("cache unit present after expiry time, wanted nothing")
	}
}

func TestCacheImpl_Set(t *testing.T) {
	c := NewCache(nil, nil)

	_ = c.Set("key", NewCacheUnit(7, nil))
	if val, ok := c.Get("key"); !ok {
		t.Errorf("cache unit not present after adding")
	} else if val != 7 {
		t.Errorf("value incorrect got %v wanted %v", val, 7)
	}

	_ = c.Set("key", NewCacheUnit(8, nil))
	if val, ok := c.Get("key"); !ok {
		t.Errorf("cache unit not present after overwriting")
	} else if val != 8 {
		t.Errorf("cache unit has not overwritten value, got %v wanted %v", val, 8)
	}
}

func TestCacheImpl_SetExpiry(t *testing.T) {
	c := NewCache(nil, nil)

	_ = c.Set("key", NewCacheUnit(7, nil))
	if val, ok := c.Get("key"); !ok {
		t.Errorf("cache unit not present after adding")
	} else if val != 7 {
		t.Errorf("value incorrect got %v wanted %v", val, 7)
	}

	_ = c.SetExpiry("key", NewCacheUnit(8, nil), time.Now().Add(time.Millisecond*250))
	if val, ok := c.Get("key"); !ok || val != 8 {
		t.Errorf("cache unit not (correctly) present immediately after overwriting, wanted %v got %v", 8, val)
	}

	time.Sleep(time.Millisecond * 300)
	if _, ok := c.Get("key"); ok {
		t.Errorf("cache unit present after expiry")
	}
}
