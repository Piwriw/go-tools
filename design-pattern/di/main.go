package main

import "fmt"

type DIInterface interface {
	Read()
	Write()
	GetName() string
}
type App struct {
	DIInterface DIInterface
}

type DIA struct {
}

func (D DIA) Read() {
	fmt.Println("DIA read...")
}

func (D DIA) Write() {
	fmt.Println("DIA Write...")
}

func (D DIA) GetName() string {
	return "DIA"
}

type DIB struct {
}

func (D DIB) Read() {
	fmt.Println("DIB read...")
}

func (D DIB) Write() {
	fmt.Println("DIB Write...")
}

func (D DIB) GetName() string {
	return "DIB"
}

// NewApp 构造函数，接受一个 Logger 类型的依赖并返回 App 实例
func NewApp(di DIInterface) *App {
	return &App{DIInterface: di}
}

/*
DI 依赖倒置 就是根据传入的实现 来创建对象
*/
func main() {
	dia := DIA{}
	dib := DIB{}
	app := NewApp(dia)
	app = NewApp(dib)

	app.DIInterface.Read()
	app.DIInterface.Write()
	fmt.Println(app.DIInterface.GetName())
}
