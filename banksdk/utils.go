package jdbsdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func hmac256(key, body []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}
