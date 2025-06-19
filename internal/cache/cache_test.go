package cache

import (
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := New(5 * time.Minute)
	defer cache.Stop()

	// Тестируем установку и получение значения
	cache.Set("test_key", "test_value")

	value, exists := cache.Get("test_key")
	if !exists {
		t.Error("Expected key to exist in cache")
	}

	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %v", value)
	}
}

func TestCache_GetNonExistentKey(t *testing.T) {
	cache := New(5 * time.Minute)
	defer cache.Stop()

	value, exists := cache.Get("non_existent_key")
	if exists {
		t.Error("Expected key not to exist in cache")
	}

	if value != nil {
		t.Errorf("Expected nil value, got %v", value)
	}
}

func TestCache_TTL(t *testing.T) {
	cache := New(100 * time.Millisecond) // Очень короткий TTL для теста
	defer cache.Stop()

	cache.Set("test_key", "test_value")

	// Сразу после установки значение должно быть доступно
	value, exists := cache.Get("test_key")
	if !exists || value != "test_value" {
		t.Error("Value should be available immediately after setting")
	}

	// Ждем истечения TTL
	time.Sleep(150 * time.Millisecond)

	// Значение должно быть недоступно
	_, exists = cache.Get("test_key")
	if exists {
		t.Error("Value should have expired")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := New(5 * time.Minute)
	defer cache.Stop()

	cache.Set("test_key", "test_value")
	cache.Delete("test_key")

	_, exists := cache.Get("test_key")
	if exists {
		t.Error("Key should have been deleted")
	}
}

func TestCache_Clear(t *testing.T) {
	cache := New(5 * time.Minute)
	defer cache.Stop()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	cache.Clear()

	_, exists1 := cache.Get("key1")
	_, exists2 := cache.Get("key2")

	if exists1 || exists2 {
		t.Error("All keys should have been cleared")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := New(5 * time.Minute)
	defer cache.Stop()

	// Тестируем concurrent access
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set("key", i)
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			cache.Get("key")
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			cache.Delete("key")
		}
	}()

	// Даем время горутинам выполниться
	time.Sleep(100 * time.Millisecond)

	// Если тест не упал с race condition, значит все OK
}
