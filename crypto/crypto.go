package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	//"crypto/cipher"
)

func Encrypt(data []byte, key string) []byte {
	cipher, _ := aes.NewCipher([]byte(key))
	cipher.Encrypt(data, PKCS5Pad(data))
	return data
}

func PKCS5Pad(data []byte) []byte {
	var blocksize uint8 = 16
	var padCount uint8 = 0
	padCount = blocksize - uint8(len(data)%int(blocksize))
	b := []byte{padCount}

	ba := bytes.Repeat(b, int(padCount))

	return append(data, ba...)
}

func Sha256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr
}
