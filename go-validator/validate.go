package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	IP string `validate:"hostname"`
}

func main() {
	validate := validator.New()
	v := Validator{IP: "www.baidu.com"}
	fmt.Println(validate.Struct(v))
}
