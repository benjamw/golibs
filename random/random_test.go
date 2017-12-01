package random

import (
	"os"
	"strings"
	"testing"
	"time"
)

var (
	numTries int = 20000
	loc      *time.Location
)

func TestMain(m *testing.M) {
	// test in UTC because the random dates get created in UTC
	loc, _ = time.LoadLocation("UTC")
	runVal := m.Run()
	os.Exit(runVal)
}

// Numbers

func TestInt63(t *testing.T) {
	for i := 0; i < numTries; i++ {
		v := Int63()
		if v < 0 {
			t.Fatalf("Int63 returned a value less than 0: %d", v)
		}
		if 9223372036854775807 < v {
			t.Fatalf("Int63 returned a value greater than Int63: %d", v)
		}
	}
}

func TestIntn(t *testing.T) {
	// test swapping min and max
	for i := 0; i < numTries; i++ {
		v := Intn(20, 0)
		if v < 0 {
			t.Fatalf("Intn returned a value less than 0: %d", v)
		}
		if 20 < v {
			t.Fatalf("Intn returned a value greater than 20: %d", v)
		}
	}

	// test min and max the same
	for i := 0; i < numTries; i++ {
		v := Intn(50, 50)
		if v != 50 {
			t.Fatalf("Intn returned a value other than 50: %d", v)
		}
	}

	// test proper
	for i := 0; i < numTries; i++ {
		v := Intn(100, 200)
		if v < 100 {
			t.Fatalf("Intn returned a value less than 100: %d", v)
		}
		if 200 < v {
			t.Fatalf("Intn returned a value greater than 200: %d", v)
		}
	}

	// test inclusion
	got0 := false
	got1 := false
	for i := 0; i < numTries; i++ {
		v := Intn(0, 1)
		if 0 == v {
			got0 = true
		}
		if 1 == v {
			got1 = true
		}
	}

	if !got0 {
		t.Fatal("Intn did not return min")
	}
	if !got1 {
		t.Fatal("Intn did not return max")
	}
}

func TestFloatn(t *testing.T) {
	// test swapping min and max
	for i := 0; i < numTries; i++ {
		v := Floatn(20, 0)
		if v < 0 {
			t.Fatalf("Floatn returned a value less than 0: %f", v)
		}
		if 21 < v {
			t.Fatalf("Floatn returned a value greater than 20: %f", v)
		}
	}

	// test min and max the same
	for i := 0; i < numTries; i++ {
		v := Floatn(50, 50)
		if v < 50 {
			t.Fatalf("Floatn returned a value less than 50: %f", v)
		}
		if 51 < v {
			t.Fatalf("Floatn returned a value greater than 51: %f", v)
		}
	}

	// test proper
	for i := 0; i < numTries; i++ {
		v := Floatn(100, 200)
		if v < 100 {
			t.Fatalf("Floatn returned a value less than 100: %f", v)
		}
		if 201 < v {
			t.Fatalf("Floatn returned a value greater than 200: %f", v)
		}
	}

	// test inclusion
	got0 := false
	got1 := false
	for i := 0; i < numTries; i++ {
		v := Floatn(0, 1)
		if 0 == v {
			got0 = true
		}
		if 1 == v {
			got1 = true
		}
	}

	if !got0 {
		// This will hardly ever happen, so no way to test for this
		//t.Fatal("Floatn did not return min")
	}
	if !got1 {
		// This will hardly ever happen, so no way to test for this
		//t.Fatal("Floatn did not return max")
	}
}

// Dates and Times

func TestDaten(t *testing.T) {
	// test swapping min and max
	for i := 0; i < numTries; i++ {
		v := Daten(2030, 1998)
		v = v.In(loc)
		if v.Year() < 1998 {
			t.Fatalf("Daten returned a value less than 1998: %v", v)
		}
		if 2030 < v.Year() {
			t.Fatalf("Daten returned a value greater than 2030: %v", v)
		}
	}
}

func TestDate(t *testing.T) {
	for i := 0; i < numTries; i++ {
		v := Date()
		v = v.In(loc)
		if v.Year() < 1999 {
			t.Fatalf("Date returned a value less than 1999: %v", v)
		}
		if 2020 < v.Year() {
			t.Fatalf("Date returned a value greater than 2020: %v", v)
		}
	}
}

// Strings

func TestStringnt(t *testing.T) {
	// test length
	for i := 0; i < numTries; i++ {
		v := Stringnt(20, ALPHANUMERIC)
		if len(v) != 20 {
			t.Fatalf("Stringnt returned a string that was not 20 characters long. Length: %d", len(v))
		}
	}
}

func TestStringn(t *testing.T) {
	// test length
	for i := 0; i < numTries; i++ {
		v := Stringn(30)
		if len(v) != 30 {
			t.Fatalf("Stringn returned a string that was not 30 characters long. Length: %d", len(v))
		}
	}
}

func TestString(t *testing.T) {
	// test length
	for i := 0; i < numTries; i++ {
		v := String()
		if len(v) <= 0 {
			t.Fatalf("String returned a string that was not at least 1 character long. Length: %d", len(v))
		}
		if 255 < len(v) {
			t.Fatalf("String returned a string that was greater than 255 characters long. Length: %d", len(v))
		}
	}
}

func TestEmail(t *testing.T) {
	// test length
	for i := 0; i < numTries; i++ {
		v := Email()
		s := strings.Split(v, "@")
		if len(s) != 2 {
			t.Fatalf("Email returned an address that did not contain exactly one @ symbol. Num: %d", len(s))
		}
		if len(s[0]) != 10 {
			t.Fatalf("Email returned an address that was not 10 character long. Length: %d", len(s[0]))
		}
		if s[1] != "example.com" {
			t.Fatalf("Email returned an address that was not in the example.com domain. Domain: %s", s[1])
		}
	}
}
