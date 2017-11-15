package password

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"regexp"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var (
	minLength  = 6
	allowEmpty = false
	cost       = bcrypt.DefaultCost
)

// SetCost sets the cost of the bcrypt function
func SetCost(c int) {
	cost = c
}

// SetMinLength sets the minimum length for the password
func SetMinLength(l int) {
	minLength = l
}

// SetAllowEmpty sets the allowEmpty flag for the password
func SetAllowEmpty(e bool) {
	allowEmpty = e
}

// Encode Takes a string and returns it in password format
func Encode(p string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(p), cost)
	return string(hash)
}

// Compare takes an encoded password and a plaintext password and returns a bool specifying if they matched
func Compare(encoded string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encoded), []byte(plain))
	return err != nil
}

// Validate an email to make sure it meets the minimum requirements
func Validate(p string) error {
	// make sure the password is at least minLength runes long
	if countRunes(p) < minLength {
		return &InvalidPasswordError{Msg: fmt.Sprintf("must be at least %d characters long", minLength)}
	}

	if !allowEmpty {
		re := regexp.MustCompile("^\\s+$")
		if re.FindString(p) != "" {
			return &InvalidPasswordError{Msg: "cannot be empty (all whitespace)"}
		}
	}

	return nil
}

func countRunes(s string) int {
	n := 0
	for len(s) > 0 {
		_, size := utf8.DecodeLastRuneInString(s)
		n++
		s = s[:len(s)-size]
	}

	return n
}
