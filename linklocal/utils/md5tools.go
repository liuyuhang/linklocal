package utils

import (
	"crypto/md5"
	"fmt"
)



func GetMd5s(s string) string {
	salt := "The world is his who enjoys it."
	h := md5.New()
	h.Write([]byte(s))
	h.Write([]byte(salt))
	return fmt.Sprintf("%x", h.Sum(nil))
}
