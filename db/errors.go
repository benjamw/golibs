package db

import (
	"fmt"
	"net/http"
)

// UnfoundObjectError gets thrown when an object is not found in the database
type UnfoundObjectError struct {
	EntityType string // model.EntityType() response ("Vendor", "Asset", etc)
	Key        string // the entities search key
	Value      string // the entities search value
	Err        error  // the original error
}

func (e *UnfoundObjectError) Error() string {
	return fmt.Sprintf("Cannot find %s with %s of '%s'", e.EntityType, e.Key, e.Value)
}

func (e *UnfoundObjectError) Code() int {
	return http.StatusNotFound
}

// MissingKeyError gets thrown when the datastore key has not been attached to an object before calling PostLoad
type MissingKeyError struct {
}

func (e *MissingKeyError) Error() string {
	return "Key not attached"
}

func (e *MissingKeyError) Code() int {
	return http.StatusServiceUnavailable
}

// NotARealDelete gets thrown when PreDelete hijacks the delete process
type NotARealDelete struct {
}

func (e *NotARealDelete) Error() string {
	return "not really deleted"
}

func (e *NotARealDelete) Code() int {
	return 0
}

// MissingParentKeyError gets thrown when the parent key has not been attached to an object before calling PreSave
type MissingParentKeyError struct {
}

func (e *MissingParentKeyError) Error() string {
	return "missing parent key"
}

// No Code() method for MissingParentKeyError because it should not propagate to the user

// MissingRequiredError gets thrown during PreSave when an object is getting saved with missing data
type MissingRequiredError struct {
	Property string // the property that is missing
}

func (e *MissingRequiredError) Error() string {
	return fmt.Sprintf("required %s is missing", e.Property)
}

// No Code() method for MissingRequiredError because it should not propagate to the user
