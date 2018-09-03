package saves

import "io"

type Store interface {
	Open(game string) (io.ReadCloser, error)
	Create(game string) (io.WriteCloser, error)
}
