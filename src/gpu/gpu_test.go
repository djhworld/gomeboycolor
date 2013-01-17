package gpu

import "testing"

func TestReset(t *testing.T) {
	var g *GPU = new(GPU)
	g.Reset()
}
