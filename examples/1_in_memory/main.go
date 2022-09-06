// Persist and retrieve a User object using the in-memory store.
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mplewis/ez3"

	// Your code must import the driver for your S3 backend. In this example, we use the in-memory storage.
	// To use AWS S3, you would instead import "gocloud.dev/blob/s3blob" and use an "s3://" URL.
	// Available drivers and more info: https://gocloud.dev/howto/blob/
	_ "gocloud.dev/blob/memblob"
)

// For example use only: Panic if an unexpected error occurs.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

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
	store, err := ez3.New(context.Background(), "mem://")
	check(err)

	// Create a new User and store it as `my-user`
	u := User{Name: "John", Email: "john@gmail.com"}
	err = store.Set("my-user", &u)
	check(err)
	fmt.Printf("Stored user: %+v\n", u)

	// List all keys starting with `u`
	keys, err := store.ListAll("u")
	check(err)
	fmt.Printf("Found prefixed keys: %v\n", keys)

	// Build a new User struct from the stored `my-user` data
	var u2 User
	err = store.Get("my-user", &u2)
	check(err)
	fmt.Printf("Retrieved user: %+v\n", u2)

	// Delete the `my-user` key
	err = store.Del("my-user")
	check(err)
	fmt.Println("Deleted user")

	// Fail to fetch `my-user` after deletion
	var u3 User
	err = store.Get("my-user", &u3)
	fmt.Printf("Attempted retrieval of user: %v\n", err)
}
