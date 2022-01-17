// Persist and retrieve a User object using Backblaze B2
package main

import (
	"encoding/json"
	"fmt"
	"github.com/mplewis/ez3"
	"log"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u User) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) Deserialize(data []byte) error {
	return json.Unmarshal(data, &u)
}

func main() {
	// Configure a custom S3 client to use Backblaze B2 rather than AWS S3.
	// You'll need to get the endpoint and region from your bucket's `endpoint` field
	// in the Backblaze console: https://secure.backblaze.com/b2_buckets.htm
	client, err := ez3.NewS3Client(ez3.S3ClientArgs{
		Endpoint: "https://s3.us-west-001.backblazeb2.com",
		Region:   "us-west-001",
	})
	if err != nil {
		log.Panic(err)
	}

	// Create the Backblaze B2-backed store using your environment's credentials
	// (env vars AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY).
	store, err := ez3.NewS3(ez3.S3Args{
		Client:    client,
		Bucket:    "mplewis-s3kv-test",
		Namespace: "some-directory",
	})
	if err != nil {
		log.Panic(err)
	}

	// Create a new User and store it as `user`
	u := User{Name: "John", Email: "john@gmail.com"}
	err = store.Set("user", &u)
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
