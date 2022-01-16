package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Serdeable is a serializable data type.
type Serdeable interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

// User is a user of the system.
type User struct {
	Name  string
	Email string
}

// Serialize serializes the user.
func (u User) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

// Deserialize deserializes the user.
func (u User) Deserialize(data []byte) error {
	return json.Unmarshal(data, &u)
}

func main() {
	u := User{Name: "John", Email: "john@gmail.com"}
	items := []Serdeable{u}
	for _, item := range items {
		data, err := item.Serialize()
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(string(data))
	}
}
