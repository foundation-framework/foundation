package sockets

import "reflect"

func copyInterfaceValue(i interface{}) interface{} {
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return reflect.New(reflect.ValueOf(i).Elem().Type()).Interface()
	} else {
		return reflect.New(reflect.TypeOf(i)).Interface()
	}
}
