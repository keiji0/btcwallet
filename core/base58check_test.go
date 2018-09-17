package core

import (
	"bytes"
	"testing"
)

func TestBase58Check(t *testing.T) {
	versionPrefix := byte(0x01)
	payload := []byte{0x12, 0x13, 0x14}

	base58check := NewBase58Check(versionPrefix, payload)
	if base58check.VersionPrefix != versionPrefix {
		t.Errorf("VersionPrefixが格納されていません")
	}
	if !bytes.Equal(payload, base58check.Payload) {
		t.Errorf("Payloadが格納されていません")
	}

	encodeString := base58check.String()
	base58CheckDecode, err := ImportBase58Check(encodeString)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(base58check.Bytes(), base58CheckDecode.Bytes()) {
		t.Errorf("復元したbase58checkが一致しません")
	}

	base58CheckDecode.VersionPrefix = byte(0x02)
	if bytes.Equal(base58check.Bytes(), base58CheckDecode.Bytes()) {
		t.Errorf("内容が変わっているにもかかわらず一致しました")
	}
}
