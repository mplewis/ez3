# ez3

[![Go Reference](https://pkg.go.dev/badge/github.com/mplewis/ez3.svg)](https://pkg.go.dev/github.com/mplewis/ez3)

ez3 makes it easy to use [AWS S3](https://aws.amazon.com/s3/) as a key-value store. It handles serialization automatically as long as your data structs implement the `ez3.Serializable` interface.

ez3 uses [AWS SDK for Go V2](https://aws.github.io/aws-sdk-go-v2/) under the hood, so it works with S3-compatible cloud storage providers aside from AWS, such as [Backblaze B2](examples/3_backblaze_b2/main.go).

# Usage

```go
// Connect to AWS S3 using your local credentials
store, err := ez3.NewS3(ez3.S3Args{
    Bucket:    "some-bucket",
    Namespace: "some-directory",
})
check(err)

// Persist a User struct (implements Serializable) to S3.
u1 := User{Name: "John", Email: "john@gmail.com"}
err = store.Set("user", &u1)
check(err)

// Then fetch the data back from S3 into a User struct.
var u2 User
err := store.Get("user", &u2)
check(err)
```

See the [`examples` directory](examples) for complete examples.
