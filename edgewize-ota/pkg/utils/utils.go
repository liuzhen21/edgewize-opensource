package utils

import (
	"math/rand"
	"os"
	"strings"
	"time"
)

func GetImmediateSubdirectories(dirPath string) ([]string, error) {
	var subdirectories []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subdirectories = append(subdirectories, entry.Name())
		}
	}

	return subdirectories, nil
}

// FileIsExist check file is exist
func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

// 生成指定长度的随机字符串
func RandomString(length int) string {
	// 定义随机字符串的字符集
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 设置随机数种子
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成随机字符串
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charSet[rd.Intn(len(charSet))]
	}

	return string(randomString)
}

func IP2TempEdgeName(ip string) string {
	return strings.ReplaceAll(ip, ".", "-")
}
