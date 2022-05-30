package idlib

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

//CreateMD5 CreateMD5
func CreateMD5() string {
	identifyData := make([]byte, 16)
	io.ReadFull(rand.Reader, identifyData)
	str := base64.URLEncoding.EncodeToString(identifyData)
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//Encrypt32MD5 Encrypt32MD5
func Encrypt32MD5(content string) string {
	h := md5.New()
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}
