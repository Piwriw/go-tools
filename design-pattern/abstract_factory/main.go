package main

import "fmt"

type AbstractApple interface {
	ShowApple()
}

type AbstractBanana interface {
	ShowBanana()
}

type AbstractPear interface {
	ShowPear()
}
type AbsFactory interface {
	CreateApple() AbstractApple
	CreateBanana() AbstractBanana
	CreatePear() AbstractPear
}
type AppleChina struct{}

func (*AppleChina) ShowApple() {
	fmt.Println("China Apple Show")
}

type BananaChina struct{}

func (*BananaChina) ShowBanana() {
	fmt.Println("China Banana Show")
}

type PearChina struct{}

func (*PearChina) ShowPear() {
	fmt.Println("China Pear Show")
}

/*
China  Factory
*/
type ChinaFactory struct{}

func (cf *ChinaFactory) CreateApple() AbstractApple {
	var apple AbstractApple
	apple = new(AppleChina)
	return apple
}

func (cf *ChinaFactory) CreateBanana() AbstractBanana {
	var banana AbstractBanana
	banana = new(BananaChina)
	return banana
}
func (cf *ChinaFactory) CreatePear() AbstractPear {
	var pear AbstractPear
	pear = new(PearChina)
	return pear
}

/*  日本产品族 */
type JapanPear struct{}
type JapanFactory struct{}

func (jp *JapanPear) ShowPear() {
	fmt.Println("Japan Pear")
}
func (cf *JapanFactory) CreatePear() AbstractPear {
	var pear AbstractPear

	pear = new(JapanPear)

	return pear
}

/*  美国产品族 */
type AmericanApple struct{}

func (aa *AmericanApple) ShowApple() {
	fmt.Println("美国苹果")
}

type AmericanBanana struct{}

func (ab *AmericanBanana) ShowBanana() {
	fmt.Println("美国香蕉")
}

type AmericanPear struct{}

func (ap *AmericanPear) ShowPear() {
	fmt.Println("美国梨")
}

type AmericanFactory struct{}

func (af *AmericanFactory) CreateApple() AbstractApple {
	var apple AbstractApple

	apple = new(AmericanApple)

	return apple
}

func (af *AmericanFactory) CreateBanana() AbstractBanana {
	var banana AbstractBanana

	banana = new(AmericanBanana)

	return banana
}

func (af *AmericanFactory) CreatePear() AbstractPear {
	var pear AbstractPear

	pear = new(AmericanPear)

	return pear
}

/*
	 Abstract Factory
		1. 有一个接口的抽象工厂，描述了产品的抽象
	 type AbsFactory interface {
		CreateApple() AbstractApple
		CreateBanana() AbstractBanana
		CreatePear() AbstractPear
	}
 2. 通过new 创建具体的工厂 aFac = new(AmericanFactory)
 3. 通过aFac.CreateApple() 创建真正的产品实例
 4. 通过aApple.ShowApple() 调用实例的方法
*/
// 抽象工厂模式的优缺点
//优点：
//1.  拥有工厂方法模式的优点
//2. 当一个产品族中的多个对象被设计成一起工作时，它能够保证客户端始终只使用同一个产品族中的对象。
//3   增加新的产品族很方便，无须修改已有系统，符合“开闭原则”。
//
//缺点：
//1. 增加新的产品等级结构麻烦，需要对原有系统进行较大的修改，甚至需要修改抽象层代码，这显然会带来较大的不便，违背了“开闭原则”。
//
//3.3.5 适用场景
//(1) 系统中有多于一个的产品族。而每次只使用其中某一产品族。可以通过配置文件等方式来使得用户可以动态改变产品族，也可以很方便地增加新的产品族。
//(2) 产品等级结构稳定。设计完成之后，不会向系统中增加新的产品等级结构或者删除已有的产品等级结构。
func main() {
	//需求1: 需要美国的苹果、香蕉、梨 等对象
	//1-创建一个美国工厂
	var aFac AbsFactory
	aFac = new(AmericanFactory)

	//2-生产美国苹果
	var aApple AbstractApple
	aApple = aFac.CreateApple()
	aApple.ShowApple()

	//3-生产美国香蕉
	var aBanana AbstractBanana
	aBanana = aFac.CreateBanana()
	aBanana.ShowBanana()

	//4-生产美国梨
	var aPear AbstractPear
	aPear = aFac.CreatePear()
	aPear.ShowPear()

	//需求2: 需要中国的苹果、香蕉
	//1-创建一个中国工厂
	cFac := new(ChinaFactory)

	//2-生产中国苹果
	cApple := cFac.CreateApple()
	cApple.ShowApple()

	//3-生产中国香蕉
	cBanana := cFac.CreateBanana()
	cBanana.ShowBanana()
}
