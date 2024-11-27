package utils

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"
)

func Hashsec(url string) string {
	// 将 url 转为小写
	url = strings.ToLower(url)

	// 创建 MD5 哈希实例
	m := md5.New()
	m.Write([]byte(url))
	// hash1 := m.Sum(nil)

	// 使用正则提取域名部分
	re1 := regexp.MustCompile("://([^/]*)")
	re2 := regexp.MustCompile("://([\\S\\s]*)")
	matches1 := re1.FindStringSubmatch(url)
	matches2 := re2.FindStringSubmatch(url)

	// 删除端口号部分
	s2 := matches1[1]
	s8 := matches2[1]
	rePort := regexp.MustCompile(`:\d{1,}$`)
	s9 := rePort.ReplaceAllString(s8, "")
	s2 = rePort.ReplaceAllString(s2, "")

	// 对 s2 再次生成 MD5
	m2 := md5.New()
	m2.Write([]byte(s2))
	hash2 := m2.Sum(nil)
	s3 := fmt.Sprintf("%x", hash2)

	// 截取前 22 个字符
	s4 := s3[:len(s3)-22]

	// 对 hash2 结果和 s9 进行 MD5 再次处理
	m3 := md5.New()
	m3.Write([]byte(s3))
	m3.Write([]byte(s9))
	finalHash := m3.Sum(nil)

	// 使用 s4 和 finalHash 生成最终的哈希结果
	hashResult := fmt.Sprintf("%s-%x", s4, finalHash)

	return hashResult
}
