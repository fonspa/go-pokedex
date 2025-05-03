package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cases := []struct {
		key   string
		value []byte
	}{
		{
			key:   "key1",
			value: []byte("val1"),
		},
		{
			key:   "key2",
			value: []byte("val2"),
		},
	}
	cache := NewCache(5 * time.Second)
	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			cache.Add(c.key, c.value)
			v, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find value for key %s", c.key)
			}

			if string(v) != string(c.value) {
				t.Errorf("expected value %v, got %v", c.value, v)
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 100 * time.Millisecond
	const waitTime = baseTime + 100*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
