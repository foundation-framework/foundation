package utils

import "reflect"

func CopyInterfaceValue(i interface{}) interface{} {
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return reflect.New(reflect.ValueOf(i).Elem().Type()).Interface()
	} else {
		return reflect.New(reflect.TypeOf(i)).Interface()
	}
}
