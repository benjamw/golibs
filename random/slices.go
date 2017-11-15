package random

import (
	"crypto/rand"
)

func ByteN(n int) (b []byte, err error) {
	b = make([]byte, n)
	_, err = rand.Read(b)
	return
}
