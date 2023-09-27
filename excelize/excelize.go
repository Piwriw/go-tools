package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

/*
go get github.com/xuri/excelize/v2
表格库
WEB URL https://xuri.me/excelize/zh-hans/
*/

// Create 创建 表格
func Create() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 创建一个工作表
	index, err := f.NewSheet("Sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 设置单元格的值
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs("./excelize/Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func GetExcelList() {
	f, err := excelize.OpenFile("./excelize/Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取工作表中指定单元格的值
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
func main() {
	GetExcelList()
}
