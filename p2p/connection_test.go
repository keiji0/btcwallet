package p2p

import (
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/keiji0/btcwallet/core"
	"github.com/keiji0/btcwallet/protocol"
)

func TestConnection(t *testing.T) {

	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:18333")
	if err != nil {
		log.Fatalln(err)
	}
	conn := NewConnection(addr, core.MainNetwork)
	if err := conn.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	{
		toAddr := protocol.NetAddress{
			IP:   addr.IP,
			Port: protocol.NetPort(addr.Port),
		}
		fromAddr := protocol.NetAddress{
			IP:   addr.IP,
			Port: protocol.NetPort(addr.Port),
		}
		msg := protocol.NewMsgVersion(&toAddr, &fromAddr, 0, 0)

		if err := conn.Send(msg); err != nil {
			log.Fatalln(err)
		}
	}

	{
		msg, err := conn.Receive()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%v\n", msg)
	}
}
