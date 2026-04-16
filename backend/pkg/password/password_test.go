package password

import (
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// Feature: itsm-web-app, Property 11: Password Hashing Non-Reversibility
// **Validates: Requirements 3.8**
//
// For any user password stored in the database, the stored value SHALL be a
// bcrypt hash (not plaintext), and verifying the original password against the
// stored hash SHALL return true while any other string SHALL return false.

// passwordGen generates non-empty printable passwords (1–72 bytes, bcrypt's max input).
func passwordGen() *rapid.Generator[string] {
	return rapid.StringMatching(`[\x20-\x7E]{1,72}`)
}

// TestProperty_PasswordHashNonReversibility_OriginalVerifies checks that
// Hash(pw) always produces a hash that Verify(pw, hash) accepts.
func TestProperty_PasswordHashNonReversibility_OriginalVerifies(t *testing.T) {
	// Feature: itsm-web-app, Property 11: Password Hashing Non-Reversibility
	rapid.Check(t, func(t *rapid.T) {
		pw := passwordGen().Draw(t, "password")

		hash, err := Hash(pw)
		if err != nil {
			t.Fatalf("Hash(%q) returned error: %v", pw, err)
		}

		if !Verify(pw, hash) {
			t.Fatalf("Verify(original, hash) returned false for password %q", pw)
		}
	})
}

// TestProperty_PasswordHashNonReversibility_DifferentPasswordRejected checks
// that for any two distinct passwords, Hash(pw1) does NOT verify against pw2.
func TestProperty_PasswordHashNonReversibility_DifferentPasswordRejected(t *testing.T) {
	// Feature: itsm-web-app, Property 11: Password Hashing Non-Reversibility
	rapid.Check(t, func(t *rapid.T) {
		pw1 := passwordGen().Draw(t, "password1")
		pw2 := passwordGen().Draw(t, "password2")

		// Ensure the two passwords are actually different.
		if pw1 == pw2 {
			t.Skip("generated identical passwords, skipping")
		}

		hash, err := Hash(pw1)
		if err != nil {
			t.Fatalf("Hash(%q) returned error: %v", pw1, err)
		}

		if Verify(pw2, hash) {
			t.Fatalf("Verify(different, hash) returned true: pw1=%q pw2=%q", pw1, pw2)
		}
	})
}

// TestProperty_PasswordHashNonReversibility_HashIsNotPlaintext checks that
// the stored hash is a bcrypt string (starts with "$2a$" or "$2b$") and
// never equals the plaintext password.
func TestProperty_PasswordHashNonReversibility_HashIsNotPlaintext(t *testing.T) {
	// Feature: itsm-web-app, Property 11: Password Hashing Non-Reversibility
	rapid.Check(t, func(t *rapid.T) {
		pw := passwordGen().Draw(t, "password")

		hash, err := Hash(pw)
		if err != nil {
			t.Fatalf("Hash(%q) returned error: %v", pw, err)
		}

		// Hash must never equal the plaintext.
		if hash == pw {
			t.Fatalf("hash equals plaintext for password %q", pw)
		}

		// bcrypt hashes start with "$2a$" or "$2b$".
		if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
			t.Fatalf("hash does not look like bcrypt (prefix check failed): %q", hash)
		}
	})
}
