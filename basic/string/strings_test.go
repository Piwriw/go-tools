package main

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"testing"
	"unicode"
)

func TestCompare(t *testing.T) {
	a := "gopher"
	b := "hello world"
	fmt.Println(strings.Compare(a, b))
	fmt.Println(strings.Compare(a, a))
	fmt.Println(strings.Compare(b, a))

	// EqualFold 忽略大小写
	fmt.Println(strings.EqualFold("GO", "go"))
	fmt.Println(strings.EqualFold("壹", "一"))
}

func TestContain(t *testing.T) {
	strs := "This is the go world"
	// 子串 substr 在 s 中，返回 true
	contains := strings.Contains(strs, "is")
	t.Log(contains)
	// chars 中任何一个 Unicode 代码点在 s 中，返回 true
	strings.ContainsAny(strs, "t")

}

/*
TestCount string的计数
*/
func TestCount(t *testing.T) {
	count1 := strings.Count("abcdefghhjjj", "h")
	count2 := strings.Count("abcdefghhjjj", "a")
	count3 := strings.Count("abcdefghhjjj", "j")
	count4 := strings.Count("abcdefghhjjj", "k")
	t.Logf("cnt1:%d,cnt2:%d,cnt3:%d,cnt4:%d", count1, count2, count3, count4)
}

func TestFiles(t *testing.T) {
	fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   "))
	// FieldsFunc 自定定义实现分割
	strings.FieldsFunc("  foo bar  baz   ", unicode.IsSpace)
}

/*
TestSpilt 分割字符串
*/
func TestSpilt(t *testing.T) {
	t.Logf("%q\n", strings.Split("a,b,c", ","))
}

/*
TestPrefix 以xx开头字符串
*/
func TestPrefix(t *testing.T) {
	t.Log(strings.HasPrefix("Gopher", "Go"))
	t.Log(strings.HasPrefix("Gopher", "C"))
}

/*
TestSuffix 以xx结尾字符串
*/
func TestSuffix(t *testing.T) {
	t.Log(strings.HasSuffix("Gopher", "er"))
	t.Log(strings.HasSuffix("Gopher", "h"))
}

/*
TestIndex 获取某段字符的index
*/
func TestIndex(t *testing.T) {
	t.Log(strings.Index("Gopher is the best coder", "is"))
	t.Log(strings.Index("Gopher is the best coder", "the"))
	t.Log(strings.LastIndex("Gopher is the best  the coder", "the"))
}

/*
TestJoin 字符串拼接
*/
func TestJoin(t *testing.T) {
	t.Log(strings.Join([]string{"name:joohwan", "age=18"}, "&"))
}

/*
TestRepeat 重复字符串
*/
func TestRepeat(t *testing.T) {
	t.Log("ba" + strings.Repeat("na", 2))
}

/*
TestMap 重写一个Map 映射
*/
func TestMap(t *testing.T) {
	mapping := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z': // 大写字母转小写
			return r + 32
		case r >= 'a' && r <= 'z': // 小写字母不处理
			return r
		case unicode.Is(unicode.Han, r): // 汉字换行
			return '\n'
		}
		return -1 // 过滤所有非字母、汉字的字符
	}
	t.Log(strings.Map(mapping, "Hello你#￥%……\n（'World\n,好Hello^(&(*界gopher..."))
}

/*
TestReplace 字符串替换
*/
func TestReplace(t *testing.T) {
	t.Log(strings.Replace("Go is the best,Go Go Go", "Go", "Java", 2))
	t.Log(strings.ReplaceAll("Go is the best,Go Go Go", "Go", "Java"))
}

/*
TestUpOrLower 大小写转化
*/
func TestUpOrLower(t *testing.T) {
	t.Log(strings.ToLower("HELLO WORLD"))
	t.Log(strings.ToUpper("hello world"))
}

/*
TestTitle
其中 Title 会将 s 每个单词的首字母大写，不处理该单词的后续字符。
ToTitle 将 s 的每个字母大写
ToTitleSpecial 将 s 的每个字母大写，并且会将一些特殊字母转换为其对应的特殊大写字母。
*/
func TestTitle(t *testing.T) {
	// new way
	c := cases.Title(language.Und, cases.NoLower)
	t.Log(c.String("this is a title "))
}

/*
// 将 s 左侧和右侧中匹配 cutset 中的任一字符的字符去掉
func Trim(s string, cutset string) string
// 将 s 左侧的匹配 cutset 中的任一字符的字符去掉
func TrimLeft(s string, cutset string) string
// 将 s 右侧的匹配 cutset 中的任一字符的字符去掉
func TrimRight(s string, cutset string) string
// 如果 s 的前缀为 prefix 则返回去掉前缀后的 string , 否则 s 没有变化。
func TrimPrefix(s, prefix string) string
// 如果 s 的后缀为 suffix 则返回去掉后缀后的 string , 否则 s 没有变化。
func TrimSuffix(s, suffix string) string
// 将 s 左侧和右侧的间隔符去掉。常见间隔符包括：'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL)
func TrimSpace(s string) string
// 将 s 左侧和右侧的匹配 f 的字符去掉
func TrimFunc(s string, f func(rune) bool) string
// 将 s 左侧的匹配 f 的字符去掉
func TrimLeftFunc(s string, f func(rune) bool) string
// 将 s 右侧的匹配 f 的字符去掉
func TrimRightFunc(s string, f func(rune) bool) string
*/
func TestFixString(t *testing.T) {
}
