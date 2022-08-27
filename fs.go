package ez3

import (
	"os"
	"path/filepath"
)

// FSEZ3 is a filesystem-based implementation of the EZ3 API.
type FSEZ3 struct {
	Path string
}

// Get retrieves a value from the filesystem.
func (e FSEZ3) Get(key string, dst Serializable) error {
	data, err := os.ReadFile(filepath.Join(e.Path, key))
	if err != nil {
		if os.IsNotExist(err) {
			return KeyNotFound
		}
		return err
	}
	return dst.Deserialize(data)
}

// Set stores a value in the filesystem.
func (e FSEZ3) Set(key string, val Serializable) error {
	data, err := val.Serialize()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(e.Path, key), data, 0644)
}

// Del removes a value from the filesystem.
func (e FSEZ3) Del(key string) error {
	return os.Remove(filepath.Join(e.Path, key))
}

// List lists all keys in the filesystem with the given prefix.
func (e FSEZ3) List(prefix string) ([]string, error) {
	var keys []string
	err := filepath.Walk(e.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		keys = append(keys, path)
		return nil
	})
	return keys, err
}

// NewFS creates a new filesystem-based EZ3 client.
func NewFS(path string) FSEZ3 {
	return FSEZ3{Path: path}
}
