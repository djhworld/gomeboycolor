// +build wasm

package saves

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"log"

	"github.com/djhworld/gomeboycolor/webworker"
)

var ErrNotFound error = errors.New("No save state found")

type WebStore struct {
	db map[string][]byte
}

func NewWebStore() *WebStore {
	ws := new(WebStore)
	ws.db = make(map[string][]byte)
	return ws
}

func (n *WebStore) Open(game string) (io.ReadCloser, error) {
	if v, ok := n.db[game]; !ok {
		return &bytesReadeCloser{bytes.NewReader([]byte{})}, ErrNotFound
	} else {
		log.Println("Found save state for", game)
		return &bytesReadeCloser{bytes.NewReader(v)}, nil
	}
}

func (n *WebStore) Create(game string) (io.WriteCloser, error) {
	return newWebWriter(game), nil
}

// PutSave puts the save data into a local map
func (n *WebStore) PutSave(game, base64Data string) error {
	if save, err := base64.StdEncoding.DecodeString(base64Data); err != nil {
		return err
	} else {
		n.db[game] = save
		return nil
	}
}

// wrapper around bytes.Reader to support a no-op close method
type bytesReadeCloser struct {
	*bytes.Reader
}

func (b *bytesReadeCloser) Close() error {
	// no-op
	return nil
}

type webWriter struct {
	game string
}

func newWebWriter(game string) *webWriter {
	w := new(webWriter)
	w.game = game
	return w
}

func (c *webWriter) Write(b []byte) (int, error) {
	webworker.SendSaveState(c.game, base64.StdEncoding.EncodeToString(b))
	return len(b), nil
}

func (c *webWriter) Close() error {
	return nil
}
