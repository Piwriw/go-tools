package score

import (
	"errors"
)

/*
DoWork 本次考核核心方法：需要实现绩点计算
输出绩点只保留2位小数
SGU 绩点计算方案 （绩点*该课程学分+...)/总学分 目前已知当前该专业总学分为170
基本要求：当前学生绩点全部为正考成绩
进阶玩法（一些新的情况）：
1. 学生有补考成绩，不能进行重复计算，并且补考通过，计算绩点以1.0 为计算
2. 学生有重修成绩，进行去重，并且按照重修成绩正常进行计算
*/
func DoWork() float64 {
	return float64(2.2)
}

/*
考核要点：excelize 的使用
GetScore 读取表格 学生考试成绩
*/
func GetScore() ([][]string, error) {
	return [][]string{}, errors.New("Err")
}
