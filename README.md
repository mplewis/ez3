# ez3

[![Go Reference](https://pkg.go.dev/badge/github.com/mplewis/ez3.svg)](https://pkg.go.dev/github.com/mplewis/ez3)

ez3 makes it easy to use [AWS S3](https://aws.amazon.com/s3/) as a key-value
store. It handles serialization automatically as long as your data structs
implement the `ez3.Serializable` interface.

ez3 uses [Go CDK](https://gocloud.dev/) under the hood, so it works with
S3-compatible cloud storage providers aside from AWS. See the CDK website for a
list of [supported storage providers](https://gocloud.dev/howto/blob/#services).

# Usage

```go
import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mplewis/ez3"

	// Import the driver package which supports your backend.
	_ "gocloud.dev/blob/s3blob"
)

// Create your data structure.
// Use of JSON is not required. You can de/serialize using any scheme you like.
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Implement the Serializable interface.
func (u User) Serialize() ([]byte, error) {
	return json.Marshal(u)
}
func (u *User) Deserialize(data []byte) error {
	return json.Unmarshal(data, &u)
}

// Go CDK Blob supports options. Here we use a prefix so that multiple apps can share this bucket.
// See docs for detailed options: https://gocloud.dev/howto/blob/
const bucketURL = "s3://my-bucket?region=us-west-1&prefix=myapp/prod/"

func main() {
	// Set up the store
	store, err := ez3.New(context.Background(), bucketURL)
	check(err)

	// Create a new User and store it as `my-user`
	u := User{Name: "John", Email: "john@gmail.com"}
	err = store.Set("my-user", &u)
	check(err)

	// Fetch the user's data
	var u2 User
	err = store.Get("my-user", &u2)
	check(err)
	fmt.Printf("Retrieved user: %+v\n", u2)
}

```

See the [`examples` directory](examples) and [test suite](ez3_suite_test.go) for
complete examples.
