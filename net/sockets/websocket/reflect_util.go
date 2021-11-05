package websocket

import "reflect"

func isPointer(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Ptr
}
