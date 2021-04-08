package keys

import (
	"math/rand"
	"time"
)

var letter = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenAppKey() string {
	bytesKey := make([]byte, 16)
	newRand := rand.New(rand.NewSource(time.Now().UnixNano())) // rand 计算得到的随机数
	for i := 0; i < 16; i++ {
		bytesKey[i] = letter[newRand.Intn(len(letter))]
	}
	return string(bytesKey)
}

func GenSecretKey() string {
	bytesKey := make([]byte, 32)
	newRand := rand.New(rand.NewSource(time.Now().UnixNano())) // rand 计算得到的随机数
	for i := 0; i < 32; i++ {
		bytesKey[i] = letter[newRand.Intn(len(letter))]
	}
	return string(bytesKey)
}
