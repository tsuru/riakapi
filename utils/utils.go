package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
)

// GenerateUsername creates a random username based on a word
func GenerateUsername(word string) string {
	username := fmt.Sprintf("tsuru_%s", word)
	return username
}

// GeneratePassword creates a password based on a word and a hash
func GeneratePassword(word, salt string) string {
	h := sha1.New()
	io.WriteString(h, word)
	io.WriteString(h, salt)
	return fmt.Sprintf("%x", h.Sum(nil))
}
