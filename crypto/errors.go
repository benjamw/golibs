package crypto

import (
	"net/http"
)

// ExpiredChallengeError is thrown when a challenge has expired
type ExpiredChallengeError struct {
}

// Error satisfies the Error interface
func (e *ExpiredChallengeError) Error() string {
	return "Challenge Expired: Please try again"
}

// Code satisfies the maje.Error interface
func (e *ExpiredChallengeError) Code() int {
	return http.StatusBadRequest
}
