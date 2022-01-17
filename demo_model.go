package main

import "encoding/json"

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
