package crypto_test

import (
	"bytes"
	"testing"

	"pwmanager/internal/crypto"
)

const testKey = "0123456789abcdef0123456789abcdef"

func TestNew_Valid(t *testing.T) {
	if _, err := crypto.New(testKey); err != nil {
		t.Fatalf("New: %v", err)
	}
}

func TestNew_InvalidKeyLength(t *testing.T) {
	for _, k := range []string{"", "short", "0123456789012345678901234567890", "012345678901234567890123456789012"} {
		if _, err := crypto.New(k); err == nil {
			t.Errorf("expected error for len=%d", len(k))
		}
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	c, _ := crypto.New(testKey)
	for _, plain := range [][]byte{
		[]byte("hello"), []byte(""), []byte("привет"),
		bytes.Repeat([]byte("A"), 10000),
	} {
		enc, err := c.Encrypt(plain)
		if err != nil { t.Fatalf("enc: %v", err) }
		dec, err := c.Decrypt(enc)
		if err != nil { t.Fatalf("dec: %v", err) }
		if !bytes.Equal(dec, plain) {
			t.Errorf("mismatch")
		}
	}
}

func TestEncrypt_UniqueNonces(t *testing.T) {
	c, _ := crypto.New(testKey)
	e1, _ := c.Encrypt([]byte("same"))
	e2, _ := c.Encrypt([]byte("same"))
	if bytes.Equal(e1, e2) {
		t.Error("nonce reuse detected")
	}
}

func TestDecrypt_Invalid(t *testing.T) {
	c, _ := crypto.New(testKey)
	for _, data := range [][]byte{nil, {}, {1,2,3}, bytes.Repeat([]byte{0xFF}, 40)} {
		if _, err := c.Decrypt(data); err == nil {
			t.Error("expected error")
		}
	}
}

func TestDecrypt_Tampered(t *testing.T) {
	c, _ := crypto.New(testKey)
	enc, _ := c.Encrypt([]byte("secret"))
	enc[len(enc)-1] ^= 0xFF
	if _, err := c.Decrypt(enc); err == nil {
		t.Error("tamper not detected")
	}
}
