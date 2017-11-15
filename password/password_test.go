package password

import (
	"testing"
)

func TestEncode(t *testing.T) {
	es := Encode("Your Mom")

	if !Compare(es, "Your Mom") {
		t.Errorf("Same password didn't compare the same twice")
	}
}

func TestValidate(t *testing.T) {
	// empty
	if err := Validate(""); err == nil {
		t.Fatalf("Validate did not return an error with an empty password")
	}

	// short
	if err := Validate("12345"); err == nil {
		t.Fatalf("Validate did not return an error with a short password")
	}

	// valid
	if err := Validate("123456"); err != nil {
		t.Fatalf("Validate returned an error with a valid password")
	}

	// six spaces
	if err := Validate("      "); err == nil {
		t.Fatalf("Validate did not return an error with a blank password")
	}

	// space wrapped (valid)
	if err := Validate(" 1234 "); err != nil {
		t.Fatalf("Validate returned an error with a space-wrapped valid password")
	}

	// valid UTF-8 password
	if err := Validate("𐍈 😃 € ¢ $"); err != nil {
		t.Fatalf("Validate returned an error with a valid UTF-8 password")
	}

	// too short UTF-8 password
	if err := Validate("€¢$"); err == nil {
		t.Fatalf("Validate did not return an error with a short UTF-8 password")
	}

}
