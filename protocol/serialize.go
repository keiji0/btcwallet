package protocol

import (
	"encoding/binary"
	"io"
	"math"
	"net"
	"time"

	"github.com/pkg/errors"
)

// Serialize は値をプロトコルに応じたデータに変換してwに書き込みます
func Serialize(w io.Writer, i interface{}) error {
	switch v := i.(type) {
	case bool, int8, int16, int32, int64, uint8, uint16, uint32, uint64, MessageMagic, ServiceFlags, Version, [16]byte, [messageCommandSize]byte, [messageChecksumSize]byte:
		return binary.Write(w, defaultByteOrder, v)

	case Uint32Time:
		return Serialize(w, uint32(time.Time(v).Unix()))

	case Int64Time:
		return Serialize(w, int64(time.Time(v).Unix()))

	case NetPort:
		return binary.Write(w, binary.LittleEndian, v)

	case VarUint:
		return serializeVarUint(w, v)

	case string:
		return serializeString(w, v)

	case UserAgentName:
		return serializeString(w, string(v))

	case NetAddress:
		return serializeNetAddress(w, v)

	default:
		return errors.Errorf("invalid type: %T", v)
	}
	return nil
}

// Deserialize はrからローカル環境で利用できるデータに変換してvに読み込みます
func Deserialize(r io.Reader, i interface{}) error {
	switch p := i.(type) {
	case *bool, *int8, *int16, *int32,
		*int64, *uint8, *uint16, *uint32, *uint64,
		*MessageMagic, *ServiceFlags, *Version,
		*[16]byte, *[messageCommandSize]byte, *[messageChecksumSize]byte:

		if err := binary.Read(r, defaultByteOrder, p); err != nil {
			return errors.Wrapf(err, "読み込みに失敗しました: %T", p)
		}

	case *Uint32Time:
		var v uint32
		if err := Deserialize(r, &v); err != nil {
			return err
		}
		*p = Uint32Time(time.Unix(int64(v), 0))

	case *Int64Time:
		var v uint64
		if err := Deserialize(r, &v); err != nil {
			return err
		}
		*p = Int64Time(time.Unix(int64(v), 0))

	case *NetPort:
		if err := binary.Read(r, binary.LittleEndian, p); err != nil {
			return errors.Wrapf(err, "読み込みに失敗しました: %T", p)
		}

	case *VarUint:
		return deserializeVarUint(r, p)

	case *string:
		return deserializeString(r, p)

	case *UserAgentName:
		var s string
		if err := deserializeString(r, &s); err != nil {
			return err
		}
		*p = UserAgentName(s)

	case *NetAddress:
		return deserializeNetAddress(r, p)

	default:
		return errors.Errorf("invalid type: %T", i)
	}
	return nil
}

// BulkSerialize は一括でシリアライズをします
func BulkSerialize(w io.Writer, args ...interface{}) error {
	for _, i := range args {
		if err := Serialize(w, i); err != nil {
			return err
		}
	}
	return nil
}

// BulkDeserialize は一括でデシリアライズします
func BulkDeserialize(r io.Reader, args ...interface{}) error {
	for _, i := range args {
		if err := Deserialize(r, i); err != nil {
			return err
		}
	}
	return nil
}

// serializeVarUint は可変長数値をシリアライズする
// https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func serializeVarUint(w io.Writer, v VarUint) error {
	if v <= varUint8Max {
		return binary.Write(w, defaultByteOrder, byte(v))
	}
	if v <= math.MaxUint16 {
		if err := binary.Write(w, defaultByteOrder, varUint16Tag); err != nil {
			return err
		}
		return binary.Write(w, defaultByteOrder, uint16(v))
	}
	if v <= math.MaxUint32 {
		if err := binary.Write(w, defaultByteOrder, varUint32Tag); err != nil {
			return err
		}
		return binary.Write(w, defaultByteOrder, uint32(v))
	}
	if err := binary.Write(w, defaultByteOrder, varUint64Tag); err != nil {
		return err
	}
	return binary.Write(w, defaultByteOrder, uint64(v))
}

// deserializeVarUint は可変長数値をデシリアライズする
// https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func deserializeVarUint(r io.Reader, p *VarUint) error {
	var length uint8
	if err := binary.Read(r, defaultByteOrder, &length); err != nil {
		return errors.Wrap(err, "VarUintの長さの読み込みに失敗しました")
	}

	switch length {
	case varUint16Tag:
		var v uint16
		if err := binary.Read(r, defaultByteOrder, &v); err != nil {
			return errors.Wrap(err, "VarUint16の読み込みに失敗しました")
		}
		*p = VarUint(v)

	case varUint32Tag:
		var v uint32
		if err := binary.Read(r, defaultByteOrder, &v); err != nil {
			return errors.Wrap(err, "VarUint32の読み込みに失敗しました")
		}
		*p = VarUint(v)

	case varUint64Tag:
		var v uint64
		if err := binary.Read(r, defaultByteOrder, &v); err != nil {
			return errors.Wrap(err, "VarUint64の読み込みに失敗しました")
		}
		*p = VarUint(v)

	default:
		*p = VarUint(length)
	}

	return nil
}

// serializeString は可変長の文字列をシリアライズします
// https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_string
func serializeString(w io.Writer, v string) error {
	if err := serializeVarUint(w, VarUint(len(v))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(v)); err != nil {
		return errors.Wrap(err, "Stringの書き込みに失敗しました")
	}
	return nil
}

// serializeString は可変長の文字列をデシリアライズします
// https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_string
func deserializeString(r io.Reader, p *string) error {
	*p = ""
	var len VarUint
	if err := deserializeVarUint(r, &len); err != nil {
		return err
	}
	if maxStringLength < len {
		return errors.Errorf("Stringが読み込める最大長を超えました: length=%d", len)
	}
	if len == 0 {
		// 0は正常、さっさと終わらす
		return nil
	}

	buf := make([]byte, len)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return errors.Wrapf(err, "StringのBodyの読み込みに失敗しました: length=%d", len)
	}
	if VarUint(n) != len {
		return errors.Wrapf(err, "Stringの長さと読み込んだ長さが一致しません: length=%d, read_length=%d", len, n)
	}

	*p = string(buf)
	return nil
}

// serializeIP はIPAddressをシリアライズします
func serializeNetAddress(w io.Writer, v NetAddress) error {
	var ip [16]byte
	copy(ip[:], v.IP.To16())
	return BulkSerialize(w, v.Services, ip, v.Port)
}

// deserializeIP はIPAddressをデシリアライズします
func deserializeNetAddress(r io.Reader, p *NetAddress) error {
	var ip [16]byte
	if err := BulkDeserialize(r, &p.Services, &ip, &p.Port); err != nil {
		return err
	}
	p.IP = net.IP(ip[:])
	return nil
}
