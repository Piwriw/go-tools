package main

import (
	"fmt"
	"github.com/wenlng/go-captcha/captcha"
)

/*
go-captcha 验证码库
go get -u github.com/wenlng/go-captcha/captcha
*/
func main() {
	capt := captcha.GetCaptcha()

	// 生成验证码
	dots, b64, tb64, key, err := capt.Generate()
	if err != nil {
		panic(err)
		return
	}

	// 主图base64
	fmt.Println(len(b64))

	// 缩略图base64
	fmt.Println(len(tb64))

	// 唯一key
	fmt.Println(key)

	// 文本位置验证数据
	fmt.Println(dots)

}
