package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"time"
)

//RandString ...
//t=0 random 0-9
//t=1 random a-z
//t=2 random A-Z
//t=others, random all
func RandString(size int, t int) string {
	types := [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}
	result := make([]byte, size)
	all := false
	if t > 2 || t < 0 {
		all = true
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if all {
			t = rand.Intn(3)
		}
		scope, base := types[t][0], types[t][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

// Md5String ...
func Md5String(s string) string {
	ctx := md5.New()
	ctx.Write([]byte(s))
	return hex.EncodeToString(ctx.Sum(nil))
}

// StringInSlice ...
// check if a string in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ComputeHmac256 encrypt a string in hmac256
func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// AES128CBCEncrypt encrypt a string using aes128 to base64 string
func AES128CBCEncrypt(message string) (string, error) {
	aeskey := []byte("testquantalabaiptestquantalabaip")
	pass := []byte(message)
	xpass, err := AesEncrypt(pass, aeskey)
	if err != nil {
		return "", err
	}
	pass64 := base64.StdEncoding.EncodeToString(xpass)
	return pass64, nil
}

// AES128CBCDecrypt decrypt a encrypted base64 string using aes128 to string
func AES128CBCDecrypt(pass64 string) (string, error) {
	aeskey := []byte("testquantalabaiptestquantalabaip")
	bytesPass, err := base64.StdEncoding.DecodeString(pass64)
	if err != nil {
		return "", err
	}

	tpass, err := AesDecrypt(bytesPass, aeskey)
	if err != nil {
		return "", err
	}
	return string(tpass), nil
}
