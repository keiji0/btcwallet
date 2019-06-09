package protocol

import (
	"bytes"
	"net"
	"testing"
)

func TestSerializeVarUint(t *testing.T) {

	tests := []struct {
		val VarUint
		tag uint8
		len int
	}{
		{0x00, 0x00, 1},
		{0x05, 0x00, 1},
		{0xfc, 0x00, 1},
		{0xfd, varUint16Tag, 3},
		{0xff, varUint16Tag, 3},
		{0xffff, varUint16Tag, 3},
		{0x10000, varUint32Tag, 5},
		{0xffffffff, varUint32Tag, 5},
		{0x100000000, varUint64Tag, 9},
		{0xffffffffffffffff, varUint64Tag, 9},
	}
	for _, test := range tests {
		buf := &bytes.Buffer{}
		if err := serializeVarUint(buf, test.val); err != nil {
			t.Error(err)
			continue
		}

		if test.len != buf.Len() {
			t.Errorf("適切なバイト数が書き込まれていません: %d != %d", test.len, buf.Len())
			continue
		}

		if 2 <= buf.Len() {
			tag := buf.Bytes()[0]
			if test.tag != tag {
				t.Errorf("タグが一致しません: %d != %d", test.tag, tag)
				continue
			}
		}

		var varUint VarUint
		if err := deserializeVarUint(buf, &varUint); err != nil {
			t.Error(err)
			continue
		}

		if varUint != test.val {
			t.Errorf("書き込んだ値と読み込んだ値が一致しません: %d != %d", varUint, test.val)
		}

		buf.Reset()
	}
}

func TestSerializeString(t *testing.T) {
	tests := []struct {
		val string
	}{
		{""},
		{"abc"},
		{"hogehogehoge"},
	}
	for _, test := range tests {
		buf := &bytes.Buffer{}
		if err := serializeString(buf, test.val); err != nil {
			t.Error(err)
			continue
		}

		if len(test.val)+1 != buf.Len() {
			t.Errorf("適切なバイト数が書き込まれていません: %s, %d != %d", test.val, len(test.val)+1, buf.Len())
			continue
		}

		var str string
		err := deserializeString(buf, &str)
		if err != nil {
			t.Error(err)
			continue
		}

		if str != test.val {
			t.Errorf("書き込んだ値と読み込んだ値が一致しません: %s != %s", str, test.val)
		}

		buf.Reset()
	}
}

func TestSerializeNetAddress(t *testing.T) {

	tests := []struct {
		addr NetAddress
	}{
		{
			addr: NetAddress{
				Services: 0x01,
				IP:       net.IP([]byte{192, 168, 0, 1}),
				Port:     0xf0,
			},
		},
	}
	for _, test := range tests {
		buf := &bytes.Buffer{}
		if err := Serialize(buf, test.addr); err != nil {
			t.Error(err)
			continue
		}

		var addr NetAddress
		if err := Deserialize(buf, &addr); err != nil {
			t.Error(err)
			continue
		}

		if test.addr.Services != addr.Services {
			t.Errorf("Servicesが一致しません: %d != %d", test.addr.Services, addr.Services)
			continue
		}
		if bytes.Equal(test.addr.IP, addr.IP) {
			t.Errorf("IPが一致しません: %#v != %#v", test.addr.IP, addr.IP)
			continue
		}
	}
}
