package gpu

import (
	"testing"
)

func TestSprite8x8PushPop(t *testing.T) {
	sprite := NewSprite8x8()

	if sprite.IsScanlineDrawQueueEmpty() != true {
		t.Log("Expected draw queue to be empty")
		t.Fail()
	}

	sprite.PushScanlines(1, 8)

	if sprite.IsScanlineDrawQueueEmpty() != false {
		t.Log("Expected draw queue to be full of values")
		t.Fail()
	}

	for i := 0; i < 8; i++ {
		j, k := sprite.PopScanline()

		if j != i+1 {
			t.Log("Expected j =", i+1)
			t.Fail()
		}

		if k != i {
			t.Log("Expected k =", i)
			t.Fail()
		}
	}

	if sprite.IsScanlineDrawQueueEmpty() != true {
		t.Log("Expected draw queue to be empty")
		t.Fail()
	}
}

func TestSprite8x16PushPop(t *testing.T) {
	sprite := NewSprite8x16()

	if sprite.IsScanlineDrawQueueEmpty() != true {
		t.Log("Expected draw queue to be empty")
		t.Fail()
	}

	sprite.PushScanlines(1, 16)

	if sprite.IsScanlineDrawQueueEmpty() != false {
		t.Log("Expected draw queue to be full of values")
		t.Fail()
	}

	for i := 0; i < 16; i++ {
		j, k := sprite.PopScanline()

		if j != i+1 {
			t.Log("Expected j =", i+1)
			t.Fail()
		}

		if k != i {
			t.Log("Expected k =", i)
			t.Fail()
		}
	}

	if sprite.IsScanlineDrawQueueEmpty() != true {
		t.Log("Expected draw queue to be empty")
		t.Fail()
	}
}
