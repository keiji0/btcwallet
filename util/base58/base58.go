package base58

import (
	b58 "github.com/keiji0/encoding/basex"
)

var bitcoinB58, _ = b58.NewEncoder("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Encode はバイト配列をビットコイン形式のBase58に変換します
func Encode(d []byte) string {
	return bitcoinB58.Encode(d)
}

// Decode はBase58形式の文字列をバイト配列に変換します
func Decode(d string) ([]byte, error) {
	return bitcoinB58.Decode(d)
}
