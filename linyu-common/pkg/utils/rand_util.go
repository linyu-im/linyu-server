package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

const (
	accountLen = 20
	chars      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func RandomString(n int) string {
	b := make([]byte, n)
	l := big.NewInt(int64(len(chars)))
	for i := range b {
		num, err := rand.Int(rand.Reader, l)
		if err != nil {
			panic(err)
		}
		b[i] = chars[num.Int64()]
	}
	return string(b)
}

func Random6DigitCode() string {
	maxCode := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, maxCode)
	if err != nil {
		return "971120"
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func RandUsername(prefix string) string {
	return prefix + RandomString(8)
}

func GenerateAccount(accountPrefix string) string {
	var sb strings.Builder
	sb.WriteString(accountPrefix)
	// 时间戳（毫秒，36进制压缩）
	ts := time.Now().UnixMilli()
	sb.WriteString(base62(ts))
	// 随机补齐长度
	for sb.Len() < accountLen {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		sb.WriteByte(chars[n.Int64()])
	}
	return sb.String()
}

// 十进制 → base62
func base62(num int64) string {
	if num == 0 {
		return "0"
	}
	result := ""
	for num > 0 {
		result = string(chars[num%62]) + result
		num /= 62
	}
	return result
}
