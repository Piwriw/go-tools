package utils

import (
	"sync"
)

type SyncSet struct {
	m sync.Map
}

func NewSyncSet(items ...interface{}) *SyncSet {
	s := &SyncSet{}
	for _, item := range items {
		s.m.Store(item, struct{}{})
	}
	return s
}

func (s *SyncSet) Add(items ...interface{}) {
	for _, item := range items {
		s.m.Store(item, struct{}{})
	}
}

func (s *SyncSet) Contains(item interface{}) bool {
	_, ok := s.m.Load(item)
	return ok
}

func (s *SyncSet) Size() int {
	size := 0
	s.m.Range(func(_, _ interface{}) bool {
		size++
		return true
	})
	return size
}

func (s *SyncSet) Clear() {
	s.m = sync.Map{}
}

func (s *SyncSet) Equal(other *SyncSet) bool {
	if s.Size() != other.Size() {
		return false
	}

	equal := true
	s.m.Range(func(key, _ interface{}) bool {
		if !other.Contains(key) {
			equal = false
			return false // stop iteration
		}
		return true // continue iteration
	})
	return equal
}

func (s *SyncSet) IsSubset(other *SyncSet) bool {
	subset := true
	s.m.Range(func(key, _ interface{}) bool {
		if !other.Contains(key) {
			subset = false
			return false // stop iteration
		}
		return true // continue iteration
	})
	return subset
}

func (s *SyncSet) Elements() []interface{} {
	elements := make([]interface{}, 0)
	s.m.Range(func(key, _ interface{}) bool {
		elements = append(elements, key)
		return true
	})
	return elements
}

func (s *SyncSet) StringElements() []string {
	elements := make([]string, 0)
	s.m.Range(func(key, _ interface{}) bool {
		if str, ok := key.(string); ok {
			elements = append(elements, str)
		}
		return true
	})
	return elements
}

func (s *SyncSet) QueryBodyElements() []QueryBody {
	elements := make([]QueryBody, 0)
	s.m.Range(func(key, _ interface{}) bool {
		if qb, ok := key.(QueryBody); ok {
			elements = append(elements, qb)
		}
		return true
	})
	return elements
}

func (s *SyncSet) SqlIDElements() map[string][]string {
	conIDMap := make(map[string][]string)
	s.m.Range(func(key, _ interface{}) bool {
		if qb, ok := key.(QueryBody); ok {
			conIDMap[qb.ConID] = append(conIDMap[qb.ConID], qb.SqlID)
		}
		return true
	})
	return conIDMap
}
