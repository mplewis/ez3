package ez3_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mplewis/ez3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "gocloud.dev/blob/memblob"
)

func TestEz3(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EZ3 Suite")
}

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

var _ = Describe("EZ3", func() {
	It("works as expected", func() {
		store, err := ez3.New(context.Background(), "mem://")
		Expect(err).ToNot(HaveOccurred())

		u := User{Name: "John", Email: "john@gmail.com"}
		err = store.Set("user", &u)
		Expect(err).ToNot(HaveOccurred())

		keys, err := store.ListAll("u")
		Expect(err).ToNot(HaveOccurred())
		Expect(keys).To(Equal([]string{"user"}))

		var u2 User
		err = store.Get("user", &u2)
		Expect(err).ToNot(HaveOccurred())
		Expect(u2).To(Equal(u))

		err = store.Del("user")
		Expect(err).ToNot(HaveOccurred())
		keys, err = store.ListAll("u")
		Expect(err).ToNot(HaveOccurred())
		Expect(keys).To(BeEmpty())

		var u3 User
		err = store.Get("user", &u3)
		Expect(err).To(MatchError(ez3.ErrKeyNotFound))

		err = store.Del("user")
		Expect(err).To(MatchError(ez3.ErrKeyNotFound))
	})
})
