package main

import (
	"fmt"
	"strings"
)

// MemoryEZ3 is an in-memory implementation of the EZ3 API.
type MemoryEZ3 struct {
	storage map[string][]byte
}

func (e MemoryEZ3) Set(key string, val Serdeable) error {
	data, err := val.Serialize()
	if err != nil {
		return err
	}
	e.storage[key] = data
	return nil
}

func (e MemoryEZ3) Get(key string, dst Serdeable) error {
	data, ok := e.storage[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	fmt.Println(string(data))
	return dst.Deserialize(data)
}

func (e MemoryEZ3) Del(key string) error {
	delete(e.storage, key)
	return nil
}

func (e MemoryEZ3) List(prefix string) ([]string, error) {
	var keys []string
	for k := range e.storage {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

// NewMemory returns a new in-memory EZ3.
func NewMemory() EZ3 {
	return MemoryEZ3{storage: make(map[string][]byte)}
}