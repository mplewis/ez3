package main

import (
	"fmt"
	"log"
)

// Serdeable is a serializable data type.
type Serdeable interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

// EZ3 is the interface to S3.
type EZ3 interface {
	Get(key string, dst Serdeable) error
	Set(key string, val Serdeable) error
	Del(key string) error
	List(prefix string) (keys []string, err error)
}

func main() {
	//e := NewMemory()
	e, err := NewS3(S3EZ3Args{Bucket: "mplewis-s3kv-test", Namespace: "test"})
	if err != nil {
		log.Panic(err)
	}

	u := User{Name: "John", Email: "john@gmail.com"}

	err = e.Set("user", &u)
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
