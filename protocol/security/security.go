package security

import (
	"encoding/base64"
)

var SECURITY_KEY = []byte("adjawibw")

func Encrypt(origData []byte) []byte {
	if len(origData) == 0 {
		return origData
	}
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(origData)))
	base64.StdEncoding.Encode(dst, origData)
	return dst
}

func Decrypt(origData []byte) []byte {
	if len(origData) == 0 {
		return origData
	}
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(origData)))
	n, _ := base64.StdEncoding.Decode(dst, origData)
	return dst[:n]
}
