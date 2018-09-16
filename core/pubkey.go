package core

import (
	"bytes"
	"crypto/ecdsa"
)

// PublicKey は公開鍵を表します
type PublicKey struct {
	*ecdsa.PublicKey
}

// PubkeyFormat は公開鍵のフォーマットを表す型
type PubkeyFormat byte

const (
	// CompressedFormat は圧縮形式の公開鍵の先頭につけるタグ
	CompressedFormat PubkeyFormat = 0x02
	// UncompressedFormat は非圧縮形式の公開鍵の先頭につけるタグ
	UncompressedFormat PubkeyFormat = 0x04
)

// UncompressData は非圧縮形式の公開鍵データを取得します
func (p *PublicKey) UncompressData() []byte {
	return bytes.Join([][]byte{
		[]byte{byte(UncompressedFormat)},
		p.X.Bytes(),
		p.Y.Bytes(),
	}, []byte(""))
}
