// Package ez3 provides an interface to persisting structs into a key-value store such as AWS S3.
package ez3

import "errors"

// EZ3 is a persistence interface which supports de/serialization.
type EZ3 interface {
	// Get retrieves a value from the store.
	Get(key string, dst Serializable) error
	// Set stores a value in the store.
	Set(key string, val Serializable) error
	// Del removes a value from the store.
	Del(key string) error
	// List lists all keys in the store with the given prefix.
	List(prefix string) (keys []string, err error)
}

// Serializable is a data type which supports de/serialization.
// Any data stored through EZ3 must implement this interface.
type Serializable interface {
	// Serialize serializes the struct's data into bytes.
	Serialize() ([]byte, error)
	// Deserialize deserializes the given bytes into the struct's data.
	Deserialize([]byte) error
}

// KeyNotFound is the error returned when a key is not found in the store.
var KeyNotFound = errors.New("key not found in ez3 store")
