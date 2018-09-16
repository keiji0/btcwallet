package core

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	b58 "github.com/keiji0/btcwallet/util/base58"
)

func TestGenerateKey(t *testing.T) {
	// 秘密鍵を作る
	pk1, err := GeneratePrivateKey(secp256k1.S256())
	if err != nil {
		t.Error("err")
	}
	pkdata := pk1.ExportPrivateKey()
	t.Log(b58.Encode(pkdata))

	// PKをエクスポートする
	pk2, err := ImportPrivateKey(secp256k1.S256(), pkdata)
	if err != nil {
		t.Error("err")
	}
	t.Log(b58.Encode(pk2.ExportPrivateKey()))

	// 復元したPKが一致するか確認
	if !bytes.Equal(pk1.ExportPrivateKey(), pk2.ExportPrivateKey()) {
		t.Error("エクスポートキーが一致しません")
	}
}
