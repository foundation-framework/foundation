package transport

import "fmt"

func Addr(hostname, port string) string {
	return fmt.Sprintf("%s:%s", hostname, port)
}
