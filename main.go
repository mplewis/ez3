package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Serdeable is a serializable data type.
type Serdeable interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

// User is a user of the system.
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Serialize serializes the user.
func (u User) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

// Deserialize deserializes the user.
func (u *User) Deserialize(data []byte) error {
	return json.Unmarshal(data, &u)
}

// EZ3 is the interface to S3.
type EZ3 interface {
	Get(key string, dst Serdeable) error
	Set(key string, val Serdeable) error
	Del(key string) error
	List(prefix string) (keys []string, err error)
}

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

// New returns a new in-memory EZ3.
func New() EZ3 {
	return MemoryEZ3{storage: make(map[string][]byte)}
}

func main() {
	e := New()
	u := User{Name: "John", Email: "john@gmail.com"}

	err := e.Set("user", &u)
	if err != nil {
		log.Panic(err)
	}

	keys, err := e.List("u")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(keys)

	var u2 User
	err = e.Get("user", &u2)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", u2)

	err = e.Del("user")
	if err != nil {
		log.Panic(err)
	}

	var u3 User
	err = e.Get("user", &u3)
	fmt.Println(err)
}
