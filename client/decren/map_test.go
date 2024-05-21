package decren

import (
	"reflect"
	"testing"
)

func TestOrderedMap_SetAndGet(t *testing.T) {
	om := NewOrderedMap()
	om.Set("key1", "value1")
	om.Set("key2", "value2")

	if val, exists := om.Get("key1"); !exists || val != "value1" {
		t.Error("expected value1 for key1")
	}

	if val, exists := om.Get("key2"); !exists || val != "value2" {
		t.Error("expected value2 for key2")
	}

	if _, exists := om.Get("key3"); exists {
		t.Error("expected key3 to not exist")
	}
}

func TestOrderedMap_Delete(t *testing.T) {
	om := NewOrderedMap()
	om.Set("key1", "value1")
	om.Set("key2", "value2")
	om.Delete("key1")

	if _, exists := om.Get("key1"); exists {
		t.Error("expected key1 to not exist")
	}

	// Check that key2 still exists
	if _, exists := om.Get("key2"); !exists {
		t.Error("expected key2 to exist")
	}

	// Check that the size of the map is correct
	if om.Size() != 1 {
		t.Errorf("expected size to be 1, got %d", om.Size())
	}

	// Check that the index of key2 is correct
	if om.keyIndex["key2"] != 0 {
		t.Errorf("expected index of key2 to be 0, got %d", om.keyIndex["key2"])
	}

	// Check that the Keys() method returns the correct keys
	keys := om.Keys()
	if len(keys) != 1 || keys[0] != "key2" {
		t.Errorf("expected keys to be [\"key2\"], got %v", keys)
	}
}

func TestOrderedMap_Size(t *testing.T) {
	om := NewOrderedMap()
	om.Set("key1", "value1")
	om.Set("key2", "value2")

	if size := om.Size(); size != 2 {
		t.Errorf("expected size 2, got %d", size)
	}

	om.Delete("key1")

	if size := om.Size(); size != 1 {
		t.Errorf("expected size 1, got %d", size)
	}
}

func TestOrderedMap_Keys(t *testing.T) {
	om := NewOrderedMap()
	om.Set("key1", "value1")
	om.Set("key2", "value2")

	expectedKeys := []string{"key1", "key2"}
	if keys := om.Keys(); !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("expected keys %v, got %v", expectedKeys, keys)
	}

	om.Delete("key1")
	expectedKeysAfterDelete := []string{"key2"}
	if keys := om.Keys(); !reflect.DeepEqual(keys, expectedKeysAfterDelete) {
		t.Errorf("expected keys %v, got %v", expectedKeysAfterDelete, keys)
	}
}
