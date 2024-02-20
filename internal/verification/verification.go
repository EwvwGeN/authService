package verification

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	mr "math/rand"
	"net/url"
	"time"
	"unsafe"
)

var (
	key []byte
)

func GenerateVerificationCode(data string) (string, error) {
	if key == nil {
		key = generateKey(32)
	}

	plaintext := []byte(data)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptVerificationCode(data string) (string, error) {
	if len(key) == 0 {
		return "", ErrKeySession
	}
	ciphertext, _ := base64.URLEncoding.DecodeString(data)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

    return string(ciphertext), nil
}

func generateKey(n int) ([]byte) {

	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    letterIdxBits := 6
    var letterIdxMask int64 = 1<<letterIdxBits - 1
    letterIdxMax  := 63 / letterIdxBits
	src := mr.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
    for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    return *(*[]byte)(unsafe.Pointer(&b))
}

func ConvertForURL(code string) string {
	return url.QueryEscape(code)
}