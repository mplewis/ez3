package ez3

import (
	"fmt"
	"strings"
)

// MemoryEZ3 is an in-memory implementation of the EZ3 API.
type MemoryEZ3 struct {
	storage map[string][]byte
}

// Get retrieves a value from memory.
func (e MemoryEZ3) Get(key string, dst Serializable) error {
	data, ok := e.storage[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	return dst.Deserialize(data)
}

// Set stores a value in memory.
func (e MemoryEZ3) Set(key string, val Serializable) error {
	data, err := val.Serialize()
	if err != nil {
		return err
	}
	e.storage[key] = data
	return nil
}

// Del removes a value from memory.
func (e MemoryEZ3) Del(key string) error {
	delete(e.storage, key)
	return nil
}

// List lists all keys in memory with the given prefix.
func (e MemoryEZ3) List(prefix string) ([]string, error) {
	var keys []string
	for k := range e.storage {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

// NewMemory creates a new memory-based EZ3 client.
func NewMemory() EZ3 {
	return MemoryEZ3{storage: make(map[string][]byte)}
}
