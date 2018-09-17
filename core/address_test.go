package core

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	b58 "github.com/keiji0/btcwallet/util/base58"
)

func TestAddress(t *testing.T) {
	// テスト対象の秘密鍵とアドレスのペア
	// ここから生成している
	// https://jpbitcoin.com/about/whatisbitcoin3
	items := []struct {
		pk      string
		address string
	}{
		{"d41864467935fd11de1479e8712bc3df8455ace9b417db9bbb8c8622f5ba782f", "1KrGeH76a6JPQyr4DhCBzWH8GgdNNnZxNd"},
		{"80a1641bca4f685b67a802ba1a9e35d16b30f4dd054ae557803330386876d629", "1H5sZz8MXujvEPAAXmKKyC3YYhBJRYQY8"},
		{"ee9bbdf1ad6b95e616f869d4c984b7e4db2523138b61aca478fefb94543f5bc0", "1Kwz8ZycdCZWAW6uJoQhza5LJV8xv2Mmtd"},
	}

	for i, item := range items {
		rawPk, _ := hex.DecodeString(item.pk)
		pk, err := ImportPrivateKey(secp256k1.S256(), rawPk)
		if err != nil {
			t.Errorf("秘密鍵のインポートに失敗しました: %s", err)
		}

		pubKey := pk.PublicKey()
		address := NewAddress(MainNetwork, pubKey)
		addressBase58 := b58.Encode(address.Bytes())
		if addressBase58 != item.address {
			t.Errorf("No.%d アドレスが一致しません: %v, %v", i+1, addressBase58, item.address)
		}
	}
}
