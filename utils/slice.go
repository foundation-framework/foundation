package utils

import "reflect"

func InSlice(slice interface{}, value interface{}) bool {
	sliceReflect := reflect.ValueOf(slice)

	if sliceReflect.Type().Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < sliceReflect.Len(); i++ {
		if value == sliceReflect.Index(i).Interface() {
			return true
		}
	}

	return false
}

func InSliceAll(slice interface{}, values ...interface{}) bool {
	sliceReflect := reflect.ValueOf(slice)

	if sliceReflect.Type().Kind() != reflect.Slice {
		return false
	}

	sliceMap := make(map[interface{}]struct{})
	for i := 0; i < sliceReflect.Len(); i++ {
		sliceMap[sliceReflect.Index(i).Interface()] = struct{}{}
	}

	for _, value := range values {
		if _, ok := sliceMap[value]; !ok {
			return false
		}
	}

	return true
}

func InSliceAny(slice interface{}, values ...interface{}) bool {
	sliceReflect := reflect.ValueOf(slice)

	if sliceReflect.Type().Kind() != reflect.Slice {
		return false
	}

	sliceMap := make(map[interface{}]struct{})
	for i := 0; i < sliceReflect.Len(); i++ {
		sliceMap[sliceReflect.Index(i).Interface()] = struct{}{}
	}

	for _, value := range values {
		if _, ok := sliceMap[value]; ok {
			return true
		}
	}

	return false
}
