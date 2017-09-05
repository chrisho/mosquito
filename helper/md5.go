package helper

import (
	"crypto/md5"
	"encoding/hex"
)

// md5 字符
func GenerateMd5String(value string) string {
	h := md5.New()
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}
