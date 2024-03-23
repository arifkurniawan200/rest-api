package utils

import (
	"testing"
)

func TestHashAndVerifyPassword(t *testing.T) {
	password := "ArifKurniawan"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	match := VerifyPassword(password, hashedPassword)
	if !match {
		t.Errorf("VerifyPassword: expected password to match the hash, but it did not")
	}

	incorrectPassword := "KurniawanArif"
	match = VerifyPassword(incorrectPassword, hashedPassword)
	if match {
		t.Errorf("VerifyPassword: expected incorrect password to not match the hash, but it did")
	}
}
