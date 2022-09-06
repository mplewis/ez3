package ez3

import (
	"context"
	"io"

	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"
)

// Client is a persistence interface which supports de/serialization.
type Client struct {
	ctx context.Context
	b   *blob.Bucket
}

// New builds a new Client for the given bucket.
func New(ctx context.Context, bucket string) (*Client, error) {
	b, err := blob.OpenBucket(ctx, bucket)
	return &Client{ctx, b}, err
}

// Get retrieves a value from the store.
func (c *Client) Get(key string, dst Serializable) error {
	raw, err := c.b.ReadAll(c.ctx, key)
	if err != nil {
		return wrapNotFoundErr(err)
	}
	return dst.Deserialize(raw)
}

// Set sets a value in the store.
func (c *Client) Set(key string, val Serializable) error {
	raw, err := val.Serialize()
	if err != nil {
		return err
	}
	return c.b.WriteAll(c.ctx, key, raw, nil)
}

// Del deletes a value from the store.
func (c *Client) Del(key string) error {
	err := c.b.Delete(c.ctx, key)
	if err != nil {
		return wrapNotFoundErr(err)
	}
	return nil
}

// List returns an iterator for all keys in the store with the given prefix.
func (c *Client) List(prefix string) *blob.ListIterator {
	return c.b.List(&blob.ListOptions{Prefix: prefix})
}

// ListAll returns all keys in the store with the given prefix.
func (c *Client) ListAll(prefix string) (keys []string, err error) {
	iter := c.List(prefix)
	for {
		item, err := iter.Next(c.ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		keys = append(keys, item.Key)
	}
	return keys, nil
}

// wrapNotFoundErr replaces a "not found" error with EZ3.ErrKeyNotFound.
func wrapNotFoundErr(err error) error {
	if gcerrors.Code(err) == gcerrors.NotFound {
		return ErrKeyNotFound
	}
	return err
}
