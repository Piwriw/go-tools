package utils

import "errors"

var Exists = struct{}{}

var SQLIDInstanceSet = NewSet()

type Set struct {
	// struct为结构体类型的变量,更节省资源
	m map[interface{}]struct{}
}

type QueryBody struct {
	SqlID string `json:"sql_id"`
	ConID string `json:"con_id"`
}

func NewSet(items ...interface{}) *Set {
	// 获取Set的地址
	s := &Set{}
	// 声明map类型的数据结构
	s.m = make(map[interface{}]struct{})
	s.Add(items...)
	return s
}

func (s *Set) Add(items ...interface{}) error {
	for _, item := range items {
		s.m[item] = Exists
	}
	return nil
}

func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]
	return ok
}

func (s *Set) Size() int {
	return len(s.m)
}

func (s *Set) Clear() {
	s.m = make(map[interface{}]struct{})
}

func (s *Set) Equal(other *Set) bool {
	// 如果两者Size不相等，就不用比较了
	if s.Size() != other.Size() {
		return false
	}

	// 迭代查询遍历
	for key := range s.m {
		// 只要有一个不存在就返回false
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (s *Set) IsSubset(other *Set) bool {
	// s的size长于other，不用说了
	if s.Size() > other.Size() {
		return false
	}
	// 迭代遍历
	for key := range s.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (s *Set) Elements() []interface{} {
	elements := make([]interface{}, 0, len(s.m))
	for key := range s.m {
		elements = append(elements, key)
	}
	return elements
}

func (s *Set) StringElements() []string {
	elements := make([]string, 0, len(s.m))
	for key := range s.m {
		str, ok := key.(string)
		if ok {
			elements = append(elements, str)
		}

	}
	return elements
}

func (s *Set) QueryBodyElements() []QueryBody {
	elements := make([]QueryBody, 0, len(s.m))
	for key := range s.m {
		element, ok := key.(QueryBody)
		if ok {
			elements = append(elements, element)
		}

	}
	return elements
}

func (s *Set) SqlIDElements() (map[string][]string, error) {
	conIDMap := make(map[string][]string)
	for key := range s.m {
		element, ok := key.(QueryBody)
		if ok {
			conIDMap[element.ConID] = append(conIDMap[element.ConID], element.SqlID)
		} else {
			return nil, errors.New("Convert QueryBody Error")
		}
	}
	return conIDMap, nil
}
