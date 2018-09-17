package core

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
)

// PublicKey は公開鍵を表します
type PublicKey struct {
	*ecdsa.PublicKey
}

// PubkeyFormat は公開鍵のフォーマットを表す型
type PubkeyFormat byte

const (
	// CompressedEvenFormat は偶数圧縮形式の公開鍵の先頭につけるタグ
	CompressedEvenFormat PubkeyFormat = 0x02
	// CompressedOddFormat は奇数圧縮形式の公開鍵の先頭につけるタグ
	CompressedOddFormat PubkeyFormat = 0x03
	// UncompressedFormat は非圧縮形式の公開鍵の先頭につけるタグ
	UncompressedFormat PubkeyFormat = 0x04
)

// Data 圧縮の公開鍵データを取得します
func (p *PublicKey) Data() []byte {
	return bytes.Join([][]byte{
		[]byte{byte(UncompressedFormat)},
		p.X.Bytes(),
		p.Y.Bytes(),
	}, []byte(""))
}

// CompressData は圧縮形式の公開鍵データを取得します
// Xが決まるとYが定まる性質からXとYの正負情報を結合したものを圧縮公開鍵
func (p *PublicKey) CompressData() []byte {
	return bytes.Join([][]byte{
		[]byte{byte(p.compressedFormat())},
		p.X.Bytes(),
	}, []byte(""))
}

// compressedFormat は公開鍵のフォーマットタイプを返す
func (p *PublicKey) compressedFormat() PubkeyFormat {
	if isOdd(p.Y) {
		return CompressedOddFormat
	}
	return CompressedEvenFormat
}

// isOdd は奇数かどうか判定する
func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}
