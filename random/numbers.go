package random

import (
	"math/rand"
	"time"
)

var (
	s1 rand.Source
	r1 *rand.Rand
)

func init() {
	s1 = rand.NewSource(time.Now().UnixNano())
	r1 = rand.New(s1)
}

// Int31 returns a non-negative pseudo-random 31-bit integer as an int32.
func Int31() int32 {
	return r1.Int31()
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func Int63() int64 {
	return r1.Int63()
}

// Intn returns a random integer between min and max inclusive
func Intn(min int64, max int64) int64 {
	if max < min {
		max, min = min, max
	}

	diff := max - min
	if diff == 0 {
		return min
	}

	rnd := r1.Int63n(diff + 1) // +1 for inclusion

	return rnd + min
}

// Floatn returns a random float between min.0 and max.99999...
func Floatn(min int64, max int64) float32 {
	return float32(Intn(min, max)) + r1.Float32()
}
