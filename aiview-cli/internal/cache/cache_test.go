package cache

import (
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	c := New(5 * time.Second)

	c.Set("key1", "value1")

	val, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected key1 to exist")
	}
	if val != "value1" {
		t.Errorf("expected 'value1', got '%v'", val)
	}
}

func TestCache_Expiry(t *testing.T) {
	c := New(50 * time.Millisecond)

	c.Set("key1", "value1")
	time.Sleep(100 * time.Millisecond)

	_, ok := c.Get("key1")
	if ok {
		t.Error("expected key1 to be expired")
	}
}

func TestCache_Delete(t *testing.T) {
	c := New(5 * time.Second)

	c.Set("key1", "value1")
	c.Delete("key1")

	_, ok := c.Get("key1")
	if ok {
		t.Error("expected key1 to be deleted")
	}
}

func TestCache_Clear(t *testing.T) {
	c := New(5 * time.Second)

	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Clear()

	_, ok1 := c.Get("key1")
	_, ok2 := c.Get("key2")
	if ok1 || ok2 {
		t.Error("expected all keys to be cleared")
	}
}

func TestCache_Miss(t *testing.T) {
	c := New(5 * time.Second)

	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected cache miss for nonexistent key")
	}
}
