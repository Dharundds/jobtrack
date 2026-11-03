package crypto

import "testing"

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	plain := "refresh-token-value-12345"

	enc, err := EncryptToken(key, plain)
	if err != nil {
		t.Fatal(err)
	}

	out, err := DecryptToken(key, enc)
	if err != nil {
		t.Fatal(err)
	}
	if out != plain {
		t.Fatalf("got %q want %q", out, plain)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key := make([]byte, 32)
	enc, err := EncryptToken(key, "token")
	if err != nil {
		t.Fatal(err)
	}
	other := make([]byte, 32)
	other[0] = 1
	if _, err := DecryptToken(other, enc); err == nil {
		t.Fatal("expected decrypt error with wrong key")
	}
}
