package network

import (
	"fmt"
	"reflect"
	"testing"
)

var magic = []byte{0xd9, 0xb4, 0xbe, 0xf9}

func TestConnection(t *testing.T) {

	// addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:18333")
	// if err != nil {
	// 	t.Error(err)
	// }
	// conn, err := net.DialTCP("tcp4", nil, addr)
	// if err != nil {
	// 	t.Error(err)
	// }
	// defer conn.Close()

	// t.Log(conn)

	fmt.Printf("%s\n", reflect.TypeOf(8))
}
