package core

import (
	ht "github.com/keiji0/btcwallet/util/hash"
)

// Address はビットコインのアドレスを表す型
type Address struct {
	pubKey  *PublicKey
	netType NetworkType
}

// NewAddress 公開鍵からビットコインアドレスを生成
func NewAddress(netType NetworkType, pubKey *PublicKey) *Address {
	return &Address{pubKey, netType}
}

// Bytes はビットコインアドレスのバイト列を返します
// data = NetType(1) + ripem160(sh256(pubkey))
// address = data + checksum(data)
func (addr *Address) Bytes() []byte {
	versionPrefix := byte(addr.netType)
	payload := ht.Ripemd160(ht.Sha256(addr.pubKey.Data()))
	return NewBase58Check(versionPrefix, payload).Bytes()
}
