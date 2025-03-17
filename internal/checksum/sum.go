package checksum

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sum(content []byte) string {
	h := sha1.New()
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}
