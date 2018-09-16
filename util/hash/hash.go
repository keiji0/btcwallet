package hash

import (
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

// Sha256 を計算する
func Sha256(v []byte) []byte {
	ht := sha256.Sum256(v)
	return ht[:]
}

// Sha256x2 はSha256二回計算する
func Sha256x2(v []byte) []byte {
	return Sha256(Sha256(v))
}

// Ripemd160 を計算する
func Ripemd160(v []byte) []byte {
	rip := ripemd160.New()
	if _, err := rip.Write(v); err != nil {
		// rip.Wirteは失敗しないのでここにはこないはず
		return []byte{}
	}
	return rip.Sum(nil)
}
