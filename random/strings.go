package random

import (
	"strings"
)

/*
The way these values are figured is the following:

	- ...IdxBits = how many bits to represent the index of the character
		i.e.- the number of bits required to represent the number of characters
		e.g.-     1 -   2 = 1
			  3 -   4 = 2
			  5 -   8 = 3
			  9 -  16 = 4
			 17 -  32 = 5
			 33 -  64 = 6
			 65 - 128 = 7
			129 - 256 = 8
			etc...
	- ...IdxMask = 1<<...IdxBits - 1 -- this is copied as is
	- ...IdxMax = 63 / ...IdxBits    -- this is copied as is

*/

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 52 chars = 6 bits
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

const numberBytes = "0123456789"
const (
	numberIdxBits = 4 // 10 chars = 4 bits
	numberIdxMask = 1<<numberIdxBits - 1
	numberIdxMax  = 63 / numberIdxBits
)

const alphanumBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	alphanumIdxBits = 6 // 62 chars = 6 bits
	alphanumIdxMask = 1<<alphanumIdxBits - 1
	alphanumIdxMax  = 63 / alphanumIdxBits
)

const hexBytes = "0123456789ABCDEF"
const (
	hexIdxBits = 4 // 16 chars = 4 bits
	hexIdxMask = 1<<hexIdxBits - 1
	hexIdxMax  = 63 / alphanumIdxBits
)

const (
	ALPHA = iota
	NUMERIC
	ALPHANUMERIC
	HEXADECIMAL
)


// Stringnt returns a random string of length n made of kind characters
// taken mostly from http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func Stringnt(n int, kind int) string {
	bytes := letterBytes
	idxBits := uint(letterIdxBits)
	idxMask := int64(letterIdxMask)
	idxMax := letterIdxMax

	switch kind {
	case NUMERIC: // 0-9
		bytes = numberBytes
		idxBits = uint(numberIdxBits)
		idxMask = int64(numberIdxMask)
		idxMax = numberIdxMax

	case ALPHANUMERIC: // 0-9 a-z A-Z
		bytes = alphanumBytes
		idxBits = uint(alphanumIdxBits)
		idxMask = int64(alphanumIdxMask)
		idxMax = alphanumIdxMax

	case HEXADECIMAL: // 0-9 A-F
		bytes = hexBytes
		idxBits = uint(hexIdxBits)
		idxMask = int64(hexIdxMask)
		idxMax = hexIdxMax

	case ALPHA: // a-z A-Z
		fallthrough
	default:
		// do nothing
	}

	b := make([]byte, n)
	// Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, Int63(), idxMax; i >= 0; {
		if remain == 0 {
			cache, remain = Int63(), idxMax
		}
		if idx := int(cache & idxMask); idx < len(bytes) {
			b[i] = bytes[idx]
			i--
		}
		cache >>= idxBits
		remain--
	}

	return string(b)
}

// Stringn returns a random alpha string of length n
func Stringn(n int) string {
	return Stringnt(n, ALPHA)
}

// String returns a random alpha string of random length 1 < n < 255
func String() string {
	return Stringn(int(Intn(1, 255)))
}

// Emailnd returns a random email address with username of given length in the given domain
func Emailnd(n int, domain string) string {
	return strings.ToLower(Stringn(n) + "@" + domain)
}

// Emaild returns a random email address with username of length 10 in the given domain
func Emaild(domain string) string {
	return Emailnd(10, domain)
}

// Email returns a random email address with username of length 10 in the domain example.com
func Email() string {
	return Emaild("example.com")
}
