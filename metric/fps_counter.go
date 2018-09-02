package metric

const SAMPLE_SIZE = 5

type FPSCounter struct {
	samples []int
	bucket  int
}

func NewFPSCounter() *FPSCounter {
	c := new(FPSCounter)
	c.samples = make([]int, SAMPLE_SIZE, SAMPLE_SIZE)
	c.bucket = 0
	return c
}

func (f *FPSCounter) Add(sample int) {
	f.samples[f.bucket] = sample
	f.bucket++
	if f.bucket == len(f.samples) {
		f.bucket = 0
	}
}

func (f *FPSCounter) Avg() float32 {
	var average float32 = 0.0
	for _, i := range f.samples {
		average += float32(i)
	}
	result := average / float32(len(f.samples))
	return result
}
