package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

// PrivateKey アドレス生成や送金に利用するプライベートな秘密鍵になります
type PrivateKey struct {
	base *ecdsa.PrivateKey
}

// GeneratePrivateKey は秘密鍵を生成します
func GeneratePrivateKey(curve elliptic.Curve) (*PrivateKey, error) {
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{key}, nil
}

// ImportPrivateKey は秘密鍵のバイト列から秘密鍵を生成します
// pkByteにExportしたPrivateKeyのバイト列を渡します
func ImportPrivateKey(curve elliptic.Curve, pkByte []byte) (*PrivateKey, error) {
	priv := ecdsa.PrivateKey{
		D: new(big.Int).SetBytes(pkByte),
	}
	priv.PublicKey.Curve = curve
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(pkByte)
	return &PrivateKey{&priv}, nil
}

// ExportPrivateKey は生の秘密鍵をエクスポートします
func (pk *PrivateKey) ExportPrivateKey() []byte {
	return pk.base.D.Bytes()
}

// PublicKey は秘密鍵から公開鍵を取得します
func (pk *PrivateKey) PublicKey() *PublicKey {
	return &PublicKey{
		&pk.base.PublicKey,
	}
}
