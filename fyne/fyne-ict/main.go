package main

import (
	"fmt"
	"fyne-ict/gui"
	"fyne-ict/score"
)

/*

GUI中文文档 https://go-circle.cn/fyne-press/v1.0/1-getting-started/introduction.html
表格操作中文文档 https://xuri.me/excelize/zh-hans/chart/line.html
项目初始化：
1. 开启GO111MODULE=on
2. 设置中国代理 export GOPROXY=https://goproxy.cn
3. 下载依赖 go mod tidy
考核任务：
1. 难度：Easy 读取score.xlsx文件,过滤出 开课学期 课程代码 成绩 课程名称 考核方式 课程性质 课程属性 并且赋值给rows
2. 难度：基本要求Easy  进阶要求：Middle 通过获得的考试成绩数据  计算机该同学 大学平均四年绩点
3. 难度：Middle 绘制出出这个同学大学四年的课程属性的饼状图占比情况
*/
/*
printGraph  绘制出出这个同学大学四年的课程属性的饼状图占比情况
考核目的：图表的使用，以后在绘制图表的时候，代码风格类似
*/
func printGraph(scores [][]string) {
}

func main() {

	scores, err := score.GetScore()
	if err != nil {
		fmt.Printf("Reading scores.xlsx if failed , err:%s", err.Error())
		return
	}
	printGraph(scores)
	gui.App(scores)
}
