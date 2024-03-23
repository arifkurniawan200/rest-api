package utils

import (
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	MySecret := "abc&1*~#^2^#s0^=)^^7%b34"

	// Text to be encrypted and decrypted
	originalText := "Hello, World!"

	encryptedText, err := Encrypt(originalText, MySecret)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}

	decryptedText, err := Decrypt(encryptedText, MySecret)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}

	if originalText != decryptedText {
		t.Errorf("Expected decrypted text to be '%s', but got '%s'", originalText, decryptedText)
	}
}

func TestEncodeAndDecode(t *testing.T) {
	originalData := []byte("Test data")

	encodedData := Encode(originalData)

	decodedData := Decode(encodedData)

	if string(originalData) != string(decodedData) {
		t.Errorf("Expected decoded data to be '%s', but got '%s'", string(originalData), string(decodedData))
	}
}
