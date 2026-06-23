package secret

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func testKey(t *testing.T) []byte {
	t.Helper()
	key, err := newKey()
	if err != nil {
		t.Fatalf("newKey: %v", err)
	}
	return key
}

func TestSealOpenRoundTrip(t *testing.T) {
	key := testKey(t)
	plaintext := []byte("super-secret-api-key")

	blob, err := Seal(key, plaintext)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	if bytes.Contains(blob, plaintext) {
		t.Fatal("ciphertext leaks plaintext")
	}

	got, err := Open(key, blob)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("round trip mismatch: got %q want %q", got, plaintext)
	}
}

func TestSealUsesFreshNonce(t *testing.T) {
	key := testKey(t)
	a, _ := Seal(key, []byte("same"))
	b, _ := Seal(key, []byte("same"))
	if bytes.Equal(a, b) {
		t.Fatal("identical ciphertext for repeated Seal - nonce reuse")
	}
}

func TestOpenWrongKeyFails(t *testing.T) {
	blob, err := Seal(testKey(t), []byte("data"))
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	if _, err := Open(testKey(t), blob); err == nil {
		t.Fatal("expected decrypt failure with wrong key")
	}
}

func TestOpenTamperedFails(t *testing.T) {
	key := testKey(t)
	blob, _ := Seal(key, []byte("data"))
	blob[len(blob)-1] ^= 0xFF // flip a ciphertext bit
	if _, err := Open(key, blob); err == nil {
		t.Fatal("expected decrypt failure on tampered ciphertext")
	}
}

func TestMasterKeyEnvOverride(t *testing.T) {
	want := testKey(t)
	t.Setenv(envKeyName, base64.StdEncoding.EncodeToString(want))

	got, err := MasterKey(t.TempDir())
	if err != nil {
		t.Fatalf("MasterKey: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("env override key not returned")
	}
}

func TestMasterKeyEnvInvalid(t *testing.T) {
	t.Setenv(envKeyName, base64.StdEncoding.EncodeToString([]byte("too-short")))
	if _, err := MasterKey(t.TempDir()); err == nil {
		t.Fatal("expected error for wrong-length env key")
	}
}

func TestKeyFileFallbackPersists(t *testing.T) {
	dir := t.TempDir()

	first, err := fromKeyFile(dir)
	if err != nil {
		t.Fatalf("fromKeyFile (create): %v", err)
	}
	second, err := fromKeyFile(dir)
	if err != nil {
		t.Fatalf("fromKeyFile (read): %v", err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("keyfile not stable across reads")
	}
	if _, err := base64.StdEncoding.DecodeString(readFile(t, filepath.Join(dir, keyFileName))); err != nil {
		t.Fatalf("keyfile is not valid base64: %v", err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
