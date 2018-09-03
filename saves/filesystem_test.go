package saves

import (
	"io"
	"testing"
)

func TestIt(t *testing.T) {
	f := FileSystemStore{baseDir: "./"}
	w, err := f.Create("daniel.sav")
	if err != nil {
		t.FailNow()
	}
	io.WriteString(w, "foo")
}
