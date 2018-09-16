package core

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	b58 "github.com/keiji0/btcwallet/util/base58"
)

func TestAddress(t *testing.T) {
	samplePk, _ := hex.DecodeString("d41864467935fd11de1479e8712bc3df8455ace9b417db9bbb8c8622f5ba782f")
	sampleAddress := "1KrGeH76a6JPQyr4DhCBzWH8GgdNNnZxNd"

	pk, err := ImportPrivateKey(secp256k1.S256(), samplePk)
	if err != nil {
		t.Errorf("秘密鍵のインポートに失敗しました: %s", err)
	}

	pubKey := pk.PublicKey()
	address := NewAddress(MainNetwork, pubKey)
	d := b58.Encode(address.Bytes())
	if d != sampleAddress {
		t.Errorf("アドレスが一致しません: %v, %v", d, sampleAddress)
	}
	t.Log(d)
}
