package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Hmac_Sha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
