package cache

import (
	"testing"
)

func TestLRU(t *testing.T) {
	c := NewCache(
		NewLRU(
			NewLRUCache(
				5,
				nil,
				0,
			),
		),
	)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Set("d", 4)
	c.Set("e", 5)
	c.Set("f", 6)

	c.Get("d")
	c.Get("d")
	c.Get("f")
	c.Get("d")
	c.Get("f")

	for _, node := range c.GetAll().([]LRUNode) {
		t.Logf("%#v \n", node)
	}
}

func TestLFU(t *testing.T) {
	c := NewCache(
		NewLFU(
			NewLFUCache(
				5,
				nil,
				0,
			),
		),
	)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Set("d", 4)
	c.Set("e", 5)
	c.Set("f", 6)

	c.Get("d")
	c.Get("d")
	c.Get("f")
	c.Get("d")
	c.Get("f")

	for _, node := range c.GetAll().([]LFUNode) {
		t.Logf("%#v \n", node)
	}
}
