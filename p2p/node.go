package p2p

import "net"

// Node はネットワークのノードを表す型になります
type Node struct {
	IP   net.IP
	Conn *Connection
}

// NewNode はネットワークのノードを生成します
func NewNode(ip net.IP) *Node {
	return &Node{
		IP: ip,
	}
}

// func (n *Node) Connect() {
// }
