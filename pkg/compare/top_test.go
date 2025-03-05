package compare

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Total int
	Cnt   int
}
type testCaseSlice []testCase

func (s testCaseSlice) Len() int {
	return len(s)
}

func (s testCaseSlice) Less(i, j int) bool {
	return s[i].Total < s[j].Total
}

func (s testCaseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s testCaseSlice) At(i int) testCase {
	return s[i]
}

func (s testCaseSlice) Slice(start, end int) Sortable[testCase] {
	return s[start:end]
}

type IntSlice []int

func (s IntSlice) Len() int {
	return len(s)
}

func (s IntSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s IntSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s IntSlice) At(i int) int {
	return s[i]
}

func (s IntSlice) Slice(start, end int) Sortable[int] {
	return s[start:end]
}

func TestTopNSortStruct(t *testing.T) {
	// 测试空切片
	items := testCaseSlice{}
	result := TopNSort(items, 3)
	assert.Empty(t, result)

	// 测试 n 小于等于 0
	items = testCaseSlice{testCase{
		Total: 1,
		Cnt:   2,
	}, testCase{
		Total: 5,
		Cnt:   3,
	}, testCase{
		Total: 3,
		Cnt:   4,
	}}
	result = TopNSort(items, 0)
	t.Log(result)
	// expected := testCaseSlice{
	// 	{Total: 5, Cnt: 3},
	// 	{Total: 3, Cnt: 4},
	// 	{Total: 1, Cnt: 2},
	// }
	// assert.Equal(t, expected, result)
}

func TestTopNSort(t *testing.T) {
	// 测试空切片
	items := IntSlice{}
	result := TopNSort(items, 3)
	assert.Empty(t, result)

	// 测试 n 小于等于 0
	items = IntSlice{1, 2, 3, 4, 5}
	result = TopNSort(items, 0)
	assert.Empty(t, result)

	// 测试正常情况
	items = IntSlice{5, 3, 8, 1, 2}
	expected := []int{8, 5, 3}
	result = TopNSort(items, 3)
	assert.Equal(t, expected, result)

	// 测试 n 大于切片长度
	items = IntSlice{5, 3, 8, 1, 2}
	expected = []int{8, 5, 3, 2, 1}
	result = TopNSort(items, 10)
	assert.Equal(t, expected, result)
}

func TestTopN(t *testing.T) {
	// 测试空切片
	emptySlice := []int{}
	result := TopN(emptySlice, 2)
	assert.Nil(t, result)

	// 测试切片长度小于等于 n
	shortSlice := []int{1, 2}
	result = TopN(shortSlice, 3)
	assert.Equal(t, shortSlice, result)

	// 测试正常情况
	normalSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []int{1, 2, 3, 4, 5}
	result = TopN(normalSlice, 5)
	assert.Equal(t, expected, result)

	// 测试切片中有重复元素
	repeatedSlice := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	expected = []int{1, 2, 2, 3, 3, 3}
	result = TopN(repeatedSlice, 5)
	assert.Equal(t, expected, result)
}
