package helper

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"strings"
)

func Md5(sig string) string {
	newSig := md5.Sum([]byte(sig)) //转成加密编码
	// 将编码转换为字符串
	newArr := fmt.Sprintf("%x", newSig)
	//输出字符串字母都是小写，转换为大写
	return strings.ToTitle(newArr)
}

func Rand() string {
	b := make([]byte, 5)
	_, err := rand.Read(b)
	if err != nil {
		Rand()
	}
	return string(b)
}
