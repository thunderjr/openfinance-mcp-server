// Package secret resolves the master key used to encrypt the local auth
// database and provides authenticated encryption (AES-256-GCM) for the
// secret values stored at rest.
//
// The master key is resolved in the following order:
//
//  1. $OPENFINANCE_DB_KEY        - base64-encoded 32 bytes (explicit override / CI).
//  2. OS keychain                - auto-generated on first run (macOS Keychain,
//     Secret Service on Linux, Credential Manager on Windows) via go-keyring.
//  3. ~/.openfinance/master.key  - 0600 keyfile fallback when no keychain is
//     available; generated on first use.
//
// The key never lives inside the database file, so a leaked or synced DB is
// useless without access to the keychain (or the configured override).
package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
)

const (
	// KeySize is the AES-256 key length in bytes.
	KeySize = 32

	envKeyName     = "OPENFINANCE_DB_KEY"
	keyringService = "openfinance-mcp-server"
	keyringUser    = "db-encryption-key"
	keyFileName    = "master.key"
)

// MasterKey returns the 32-byte key used to seal/open secrets, resolving it
// from the environment, the OS keychain, or a keyfile under dir (in that
// order). A new key is generated and persisted on first use.
func MasterKey(dir string) ([]byte, error) {
	if raw := os.Getenv(envKeyName); raw != "" {
		key, err := decodeKey(raw)
		if err != nil {
			return nil, fmt.Errorf("secret: invalid %s: %w", envKeyName, err)
		}
		return key, nil
	}

	if key, err := fromKeyring(); err == nil {
		return key, nil
	} else if !errors.Is(err, keyring.ErrNotFound) {
		// Keychain exists but is unavailable/locked - degrade to a keyfile
		// rather than failing the whole server.
		logger.Warnf("secret: keychain unavailable (%v), falling back to keyfile", err)
		return fromKeyFile(dir)
	}

	// Keychain reachable but empty: generate, store, and return.
	key, err := newKey()
	if err != nil {
		return nil, err
	}
	if err := keyring.Set(keyringService, keyringUser, base64.StdEncoding.EncodeToString(key)); err != nil {
		logger.Warnf("secret: could not persist key to keychain (%v), falling back to keyfile", err)
		return fromKeyFile(dir)
	}
	logger.Info("secret: generated new master key in OS keychain")
	return key, nil
}

func fromKeyring() ([]byte, error) {
	raw, err := keyring.Get(keyringService, keyringUser)
	if err != nil {
		return nil, err
	}
	return decodeKey(raw)
}

func fromKeyFile(dir string) ([]byte, error) {
	path := filepath.Join(dir, keyFileName)

	raw, err := os.ReadFile(path)
	if err == nil {
		return decodeKey(string(raw))
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("secret: reading keyfile: %w", err)
	}

	key, err := newKey()
	if err != nil {
		return nil, err
	}
	encoded := base64.StdEncoding.EncodeToString(key)
	if err := os.WriteFile(path, []byte(encoded), 0o600); err != nil {
		return nil, fmt.Errorf("secret: writing keyfile: %w", err)
	}
	logger.Warnf("secret: generated master key in keyfile %s (less secure than keychain)", path)
	return key, nil
}

func newKey() ([]byte, error) {
	key := make([]byte, KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("secret: generating key: %w", err)
	}
	return key, nil
}

func decodeKey(raw string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("secret: decoding key: %w", err)
	}
	if len(key) != KeySize {
		return nil, fmt.Errorf("secret: key must be %d bytes, got %d", KeySize, len(key))
	}
	return key, nil
}

// Seal encrypts plaintext with AES-256-GCM and returns nonce||ciphertext.
func Seal(key, plaintext []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("secret: generating nonce: %w", err)
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Open decrypts a nonce||ciphertext blob produced by Seal.
func Open(key, blob []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize()
	if len(blob) < ns {
		return nil, errors.New("secret: ciphertext too short")
	}
	nonce, ciphertext := blob[:ns], blob[ns:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("secret: decrypt failed: %w", err)
	}
	return plaintext, nil
}

func newGCM(key []byte) (cipher.AEAD, error) {
	if len(key) != KeySize {
		return nil, fmt.Errorf("secret: key must be %d bytes, got %d", KeySize, len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("secret: new cipher: %w", err)
	}
	return cipher.NewGCM(block)
}
