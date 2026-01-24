package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// GenerateGravatarURL generates a Gravatar URL based on the email
func GenerateGravatarURL(email string) string {
	if email == "" {
		return "https://gravatar.com/avatar/00000000000000000000"
	}

	// Create MD5 hash of email
	h := md5.New()
	h.Write([]byte(email))
	hash := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("https://gravatar.com/avatar/%s", hash)
}
