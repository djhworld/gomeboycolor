package metric

import "testing"

func TestAvg(t *testing.T) {
	c := NewFPSCounter()
	c.Add(3.0)
	c.Add(2.0)
	c.Add(1.0)

	if c.Avg() != 1.2 {
		t.FailNow()
	}
}

func TestAddingMoreThan5Samples(t *testing.T) {
	c := NewFPSCounter()
	for i := 0; i < 20; i++ {
		c.Add(i)
	}

	if c.Avg() != 17.0 {
		t.FailNow()
	}

}
