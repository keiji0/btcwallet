package p2p

import (
	"net"

	"github.com/keiji0/btcwallet/core"
	"github.com/keiji0/btcwallet/protocol"
	"github.com/pkg/errors"
)

type connectionState int

// Connection はNodeに接続するためのデータになります
// TCPソケットを保持し、送信と受信処理を受け持ちます
type Connection struct {
	conn    *net.TCPConn
	addr    net.TCPAddr
	netType core.NetworkType
}

// NewConnection はNodeに接続するためのコネクションを生成します
func NewConnection(addr *net.TCPAddr, netType core.NetworkType) *Connection {
	c := &Connection{
		addr:    *addr,
		netType: netType,
	}
	return c
}

// Connect はノードに接続します
func (c *Connection) Connect() (err error) {
	if c.conn, err = net.DialTCP("tcp", nil, &c.addr); err != nil {
		return err
	}
	return nil
}

// Close はノードに閉じます
func (c *Connection) Close() (err error) {
	if err := c.conn.Close(); err != nil {
		return errors.Wrap(err, "コネクションを閉じるのに失敗しました")
	}
	return nil
}

// Send はコネクションに対してメッセージを送ります
func (c *Connection) Send(msg protocol.Message) error {
	if err := protocol.Send(c.conn, c.netType, msg); err != nil {
		return err
	}
	return nil
}

// Receive はコネクションからメッセージを受け取ります
func (c *Connection) Receive() (protocol.Message, error) {
	msg, err := protocol.Receive(c.conn)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
