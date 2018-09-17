package core

import (
	"bytes"

	"github.com/keiji0/btcwallet/util/base58"
	ht "github.com/keiji0/btcwallet/util/hash"
	"github.com/pkg/errors"
)

// チェックサムの長さを定義
const checksumLength = 4

// Base58Check を表す型
type Base58Check struct {
	VersionPrefix byte
	Payload       []byte
}

// NewBase58Check を生成する
func NewBase58Check(versionPrefix byte, payload []byte) *Base58Check {
	return &Base58Check{
		VersionPrefix: versionPrefix,
		Payload:       payload,
	}
}

// Bytes はBase58Checkのバイト列を取得する
func (b58c *Base58Check) Bytes() []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(b58c.VersionPrefix)
	buf.Write(b58c.Payload)
	buf.Write(calcChecksum(buf.Bytes()))
	return buf.Bytes()
}

// String はBase58Checkの文字列を取得する
func (b58c *Base58Check) String() string {
	return base58.Encode(b58c.Bytes())
}

// ImportBase58Check はBytes()からBase58Checkをインポートします
func ImportBase58Check(b58check string) (*Base58Check, error) {
	// checksum+データ部が存在しないと不正なアドレス
	if len(b58check) <= checksumLength+1 {
		return nil, errors.Errorf("チェックサム、もしくはデータ部が欠損しています: %v", b58check)
	}
	// そもそもデコードできないと不正なアドレス
	raw, err := base58.Decode(b58check)
	if err != nil {
		return nil, err
	}

	// データとチェックサムに分離して再度チェックサムを計算して一致するか判定
	versionPrefix := raw[0]
	data := raw[:len(raw)-checksumLength]
	checksum := raw[len(raw)-checksumLength:]

	if !bytes.Equal(checksum, calcChecksum(data)) {
		return nil, errors.Errorf("データのチェックサムが一致しません: %v, %s", data, checksum)
	}
	return NewBase58Check(versionPrefix, data[1:]), nil
}

// calcChecksum は受け取ったバイト列のチェックサムを計算します
func calcChecksum(dt []byte) []byte {
	hash := ht.Sha256x2(dt)
	return hash[:checksumLength]
}
