package inputoutput

import (
	"github.com/stretchrcom/testify/assert"
	"testing"
)

const (
	UP = iota
	DOWN
	LEFT
	RIGHT
	A
	B
	START
	SELECT
)

var testControlScheme ControlScheme = ControlScheme{UP, DOWN, LEFT, RIGHT, A, B, START, SELECT}

func TestWrite(t *testing.T) {
	//given
	kbh := new(KeyHandler)
	var expected byte = 0x10

	//when
	kbh.Write(0x0000, expected)

	//then
	assert.Equal(t, expected, kbh.colSelect)
}

func TestKeyDownForUp(t *testing.T) {
	var expected byte = 0x0B
	var actual byte = doKeyDownAndRead(t, ROW_1, UP)
	assert.Equal(t, expected, actual)
}

func TestKeyUpForUp(t *testing.T) {
	var expected byte = 0x0F
	var actual byte = doKeyUpAndRead(t, ROW_1, UP)
	assert.Equal(t, expected, actual)
}

func TestKeyDownForDown(t *testing.T) {
	var expected byte = 0x07
	var actual byte = doKeyDownAndRead(t, ROW_1, DOWN)
	assert.Equal(t, expected, actual)
}

func TestKeyUpForDown(t *testing.T) {
	var expected byte = 0x0F
	var actual byte = doKeyUpAndRead(t, ROW_1, DOWN)
	assert.Equal(t, expected, actual)
}

func TestKeyDownForLeft(t *testing.T) {
	var expected byte = 0x0D
	var actual byte = doKeyDownAndRead(t, ROW_1, LEFT)
	assert.Equal(t, expected, actual)
}

func TestKeyUpForLeft(t *testing.T) {
	var expected byte = 0x0F
	var actual byte = doKeyUpAndRead(t, ROW_1, LEFT)
	assert.Equal(t, expected, actual)
}

func TestKeyDownForRight(t *testing.T) {
	var expected byte = 0x0E
	var actual byte = doKeyDownAndRead(t, ROW_1, RIGHT)
	assert.Equal(t, expected, actual)
}

func TestKeyUpForRight(t *testing.T) {
	var expected byte = 0x0F
	var actual byte = doKeyUpAndRead(t, ROW_1, RIGHT)
	assert.Equal(t, expected, actual)
}

func TestKeyDownForA(t *testing.T) {
	var expected byte = 0x0E
	var actual byte = doKeyDownAndRead(t, ROW_2, A)
	assert.Equal(t, expected, actual)
}

func TestKeyUpForA(t *testing.T) {
	var expected byte = 0x0F
	var actual byte = doKeyUpAndRead(t, ROW_2, A)
	assert.Equal(t, expected, actual)
}

func TestKeyDownForB(t *testing.T) {
	var expected byte = 0x0D
	var actual byte = doKeyDownAndRead(t, ROW_2, B)
	assert.Equal(t, expected, actual)
}

func TestKeyUpForB(t *testing.T) {
	var expected byte = 0x0F
	var actual byte = doKeyUpAndRead(t, ROW_2, B)
	assert.Equal(t, expected, actual)
}

func TestKeyDownForStart(t *testing.T) {
	var expected byte = 0x07
	var actual byte = doKeyDownAndRead(t, ROW_2, START)
	assert.Equal(t, expected, actual)
}

func TestKeyUpForStart(t *testing.T) {
	var expected byte = 0x0F
	var actual byte = doKeyUpAndRead(t, ROW_2, START)
	assert.Equal(t, expected, actual)
}

func TestKeyDownForSelect(t *testing.T) {
	var expected byte = 0x0B
	var actual byte = doKeyDownAndRead(t, ROW_2, SELECT)
	assert.Equal(t, expected, actual)
}

func TestKeyUpForSelect(t *testing.T) {
	var expected byte = 0x0F
	var actual byte = doKeyUpAndRead(t, ROW_2, SELECT)
	assert.Equal(t, expected, actual)
}

//TODO: combination of keys pressed

func doKeyUpAndRead(t *testing.T, row byte, key int) byte {
	kbh := NewKeyHandler(testControlScheme)
	kbh.Write(0x0000, row)
	kbh.KeyDown(key)
	kbh.KeyUp(key)
	return kbh.Read(0x0000)
}

func doKeyDownAndRead(t *testing.T, row byte, key int) byte {
	kbh := NewKeyHandler(testControlScheme)
	kbh.Write(0x0000, row)
	kbh.KeyDown(key)
	return kbh.Read(0x0000)
}
