package network

import "net"

// Connection はNetworkに接続するためのデータになります
type Connection struct {
	conn net.TCPConn
}

// NewConnection はNetworkに接続するためのコネクションを生成します
func NewConnection(host string, port string) *Connection {
	return &Connection{}
}
