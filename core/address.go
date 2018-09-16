package core

import (
	"bytes"

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
func (addr *Address) Bytes() []byte {
	data := bytes.Join([][]byte{
		[]byte{byte(addr.netType)},
		ht.Ripemd160(ht.Sha256(addr.pubKey.UncompressData())),
	}, []byte(""))

	checksum := calcChecksum(data)
	return bytes.Join([][]byte{data, checksum}, []byte(""))
}

// calcChecksum は受け取ったバイト列のチェックサムを計算します
func calcChecksum(dt []byte) []byte {
	hash := ht.Sha256x2(dt)
	return hash[:4]
}
