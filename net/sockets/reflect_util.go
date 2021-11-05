package sockets

import "reflect"

func isPointer(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Ptr
}

func copyInterfaceValue(i interface{}) interface{} {
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return reflect.New(reflect.ValueOf(i).Elem().Type()).Interface()
	} else {
		return reflect.New(reflect.TypeOf(i)).Interface()
	}
}
