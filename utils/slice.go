package utils

import (
	"reflect"
)

type Slice struct {
	inner interface{}
}

func NewSlice(inner interface{}) *Slice {
	return &Slice{
		inner: inner,
	}
}

func (s *Slice) Has(value interface{}) bool {
	if s.inner == nil {
		return false
	}

	ref := reflect.ValueOf(s.inner)
	if ref.Type().Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < ref.Len(); i++ {
		if value == ref.Index(i).Interface() {
			return true
		}
	}

	return false
}

func (s *Slice) HasAll(values ...interface{}) bool {
	if s.inner == nil {
		return false
	}

	ref := reflect.ValueOf(s.inner)
	if ref.Type().Kind() != reflect.Slice {
		return false
	}

	valSet := make(map[interface{}]struct{})
	for i := 0; i < ref.Len(); i++ {
		valSet[ref.Index(i).Interface()] = struct{}{}
	}

	for _, value := range values {
		if _, ok := valSet[value]; !ok {
			return false
		}
	}

	return true
}

func (s *Slice) HasAny(values ...interface{}) bool {
	if s.inner == nil {
		return false
	}

	ref := reflect.ValueOf(s.inner)
	if ref.Type().Kind() != reflect.Slice {
		return false
	}

	valSet := make(map[interface{}]struct{})
	for i := 0; i < ref.Len(); i++ {
		valSet[ref.Index(i).Interface()] = struct{}{}
	}

	for _, value := range values {
		if _, ok := valSet[value]; ok {
			return true
		}
	}

	return false
}
