package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type MaxValidator struct {
	IP string `validate:"max=10"`
}

func main() {

	validate := validator.New()
	v := MaxValidator{IP: "www.baidu.com"}
	fmt.Println(validate.Struct(v))
}
