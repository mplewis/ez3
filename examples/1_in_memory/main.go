// Persist and retrieve a User object using the in-memory store.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/mplewis/ez3"
	"log"
)

// A User has a name and an email address.
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Serialize generates a representation of the User as JSON bytes.
func (u User) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

// Deserialize parses JSON bytes into a User struct.
func (u *User) Deserialize(data []byte) error {
	return json.Unmarshal(data, &u)
}

func main() {
	// Create the in-memory store
	store := ez3.NewMemory()

	// Create a new User and store it as `user`
	u := User{Name: "John", Email: "john@gmail.com"}
	err := store.Set("user", &u)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Stored user: %+v\n", u)

	// List all keys starting with `u`
	keys, err := store.List("u")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Found prefixed keys: %v\n", keys)

	// Build a new User struct from the stored `user` data
	var u2 User
	err = store.Get("user", &u2)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Retrieved user: %+v\n", u2)

	// Delete the `user` key
	err = store.Del("user")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Deleted user")

	// Fail to fetch `user` after deletion
	var u3 User
	err = store.Get("user", &u3)
	fmt.Printf("Attempted retrieval of user: %v\n", err)
}
