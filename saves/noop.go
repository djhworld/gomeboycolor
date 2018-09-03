package store

import (
	"fmt"
	"io"
)

type NoopStore struct {
	baseDir string
}

func NewNoopStore() *NoopStore {
	return new(NoopStore)
}

func (n *NoopStore) Open(game string) (io.ReadCloser, error) {
	return new(NoopContent), nil
}

func (n *NoopStore) Create(game string) (io.WriteCloser, error) {
	return new(NoopContent), nil
}

type NoopContent struct{}

func (c *NoopContent) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	return len(b), nil
}

func (c *NoopContent) Read(p []byte) (int, error) {
	return len(p), nil
}

func (c *NoopContent) Close() error {
	return nil
}
