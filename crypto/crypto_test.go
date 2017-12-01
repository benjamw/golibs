package crypto

import (
	"bytes"
	"crypto/sha1"
	"os"
	"testing"
)

var (
	txt1 = []byte("goser the gosarian")
	txt2 = []byte("stay-puft marshman")
	key1 = []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	key2 = []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x10}
)

func TestMain(m *testing.M) {
	runVal := m.Run()
	os.Exit(runVal)
}

func TestMakeSig(t *testing.T) {
	sig1 := makeSig(key1, txt1)

	if len(sig1) != sha1.Size {
		t.Fatalf("makeSig made a sig of incorrect length. Wanted: %d; Got: %d", sha1.Size, len(sig1))
	}
}

func TestEncrypt(t *testing.T) {
	clo1 := clone(txt1)

	// test that two encryptions don't match with different keys
	enc1, err := Encrypt(txt1, key1)
	if err != nil {
		t.Fatalf("Encrypt threw an error in TestEncrypt: %v", err)
	}
	if bytes.Equal(txt1, enc1) {
		t.Fatal("Encrypt did not encrypt the text")
	}

	enc2, err := Encrypt(txt1, key2)
	if err != nil {
		t.Fatalf("Encrypt threw an error in TestEncrypt: %v", err)
	}
	if bytes.Equal(txt1, enc2) {
		t.Fatal("Encrypt did not encrypt the text")
	}

	if bytes.Equal(enc1, enc2) {
		t.Fatal("Two encryptions with different keys gave the same encryption value")
	}

	// test that two encryptions don't match with different texts
	enc1, err = Encrypt(txt1, key1)
	if err != nil {
		t.Fatalf("Encrypt threw an error in TestEncrypt: %v", err)
	}
	if bytes.Equal(txt1, enc1) {
		t.Fatal("Encrypt did not encrypt the text")
	}

	enc2, err = Encrypt(txt2, key1)
	if err != nil {
		t.Fatalf("Encrypt threw an error in TestEncrypt: %v", err)
	}
	if bytes.Equal(txt2, enc2) {
		t.Fatal("Encrypt did not encrypt the text")
	}

	if bytes.Equal(enc1, enc2) {
		t.Fatal("Two encryptions with different texts gave the same encryption value")
	}

	// test that Encrypt does not mangle the original value
	if !bytes.Equal(clo1, txt1) {
		t.Fatal("Encrypt mangled the input value")
	}
}

func TestDecrypt(t *testing.T) {
	// test that encrypted text can be decrypted
	enc1, err := Encrypt(txt1, key1)
	clo1 := clone(enc1)
	if err != nil {
		t.Fatalf("Encrypt threw an error in TestDecrypt: %v", err)
	}

	dec1, err := Decrypt(clo1, key1)
	if err != nil {
		t.Fatalf("Decrypt threw an error in TestDecrypt: %v", err)
	}

	if !bytes.Equal(dec1, txt1) {
		t.Fatal("Decrypt did not return the original encryption")
	}

	// test that Decrypt does not mangle the original value
	if !bytes.Equal(clo1, enc1) {
		t.Fatal("Decrypt mangled the input value")
	}

	// test that another key will not work
	dec2, err := Decrypt(clo1, key2)
	if err != nil {
		t.Fatalf("Decrypt threw an error in TestDecrypt: %v", err)
	}

	if bytes.Equal(dec2, txt1) {
		t.Fatal("Decrypt returned the original value with the wrong key")
	}
}

func TestAddSignature(t *testing.T) {
	clo1 := clone(txt1)
	sig1 := AddSignature(txt1, key1)

	// test that AddSignature does not mangle the original value
	if !bytes.Equal(clo1, txt1) {
		t.Fatal("AddSignature mangled the input value")
	}

	// test that AddSignature is valid
	if bytes.Equal(sig1, txt1) {
		t.Fatal("AddSignature returned the same value")
	}

	if !bytes.Equal(sig1[sha1.Size:], txt1) {
		t.Fatal("AddSignature did not return the original text as part of the signature")
	}

	// test that AddSignature is different for different keys
	sig2 := AddSignature(txt1, key2)
	if bytes.Equal(sig1[:sha1.Size], sig2[:sha1.Size]) {
		t.Fatal("AddSignature returned the same signature for different keys")
	}

	// test that AddSignature is different for different values
	sig3 := AddSignature(txt2, key1)
	if bytes.Equal(sig1[:sha1.Size], sig3[:sha1.Size]) {
		t.Fatal("AddSignature returned the same signature for different values")
	}
}

func TestCheckSignature(t *testing.T) {
	sig1 := AddSignature(txt1, key1)
	clo1 := clone(sig1)

	// test that CheckSignature is valid
	out1, err := CheckSignature(sig1, key1)
	if err != nil {
		t.Fatalf("CheckSignature threw an error: %v", err)
	}

	if !bytes.Equal(out1, txt1) {
		t.Fatal("CheckSignature did not return the original value")
	}

	// test that CheckSignature does not mangle the original value
	if !bytes.Equal(clo1, sig1) {
		t.Fatal("CheckSignature mangled the input value")
	}

	// test that CheckSignature is invalid for the wrong key
	out2, err := CheckSignature(sig1, key2)
	if err == nil {
		t.Fatal("CheckSignature did not throw an error with incorrect key")
	}

	if bytes.Equal(out2, txt1) {
		t.Fatal("CheckSignature returned the original value even after an error")
	}
}

func clone(a []byte) []byte {
	b := make([]byte, len(a))
	copy(b, a)
	return b
}
