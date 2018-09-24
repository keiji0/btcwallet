package protocol

import (
	"encoding/binary"
	"net"

	"github.com/keiji0/btcwallet/core"
	"github.com/pkg/errors"
)

// Protocol Document
// https://en.bitcoin.it/wiki/Protocol_documentation

// MessageMagic はネットワークで利用するマジックタイプを表す型
type MessageMagic uint32

const (
	// MainNetMessageMagic メインネットのメッセージのマジック
	MainNetMessageMagic MessageMagic = 0xd9b4bef9
	// TestNet3MessageMagic テストネットのメッセージのマジック
	TestNet3MessageMagic MessageMagic = 0x0709110b
	// InvalidMessageMagic は不正なマジックナンバー
	InvalidMessageMagic MessageMagic = 0x00000000
)

// NetworkTypeMessageMagic はNetworkTypeからMessageMagicを取得します
func NetworkTypeMessageMagic(netType core.NetworkType) (MessageMagic, error) {
	switch netType {
	case core.MainNetwork:
		return MainNetMessageMagic, nil
	case core.TestNetwork:
		return TestNet3MessageMagic, nil
	default:
		return InvalidMessageMagic, errors.Errorf("invalid NetworkType: %v", netType)
	}
}

// Version はプロトコルのバージョンを表す型
// https://bitcoin.org/en/developer-reference#protocol-versions
type Version int32

const (
	// CurrentVersion はサポートしているプロトコルのバージョンになります
	// Bitcoin Core 0.13.2 (Jan 2017)
	CurrentVersion Version = 70015
)

// VarUint は可変長の符号なし数値を表す型
type VarUint uint64

// UserAgentName はユーザーエージェントの文字列を表す型
// エージェントは「/Name:Version/」という形式で名前を設定しネットワークへ接続する
// https://github.com/bitcoin/bips/blob/master/bip-0014.mediawiki
type UserAgentName string

// Protocolで使用するバイトオーダーを指定
// ほとんどのオーダーはリトルエンディアンだが、ポート番号はビッグエンディアンなので個別に指定する
var defaultByteOrder = binary.LittleEndian

// maxStringLength はStringのサイズの最大値
// Bitcoin Coreはメッセージのフィールドごとに決めているここでは決め打ちしておく
const maxStringLength = 0xffff

// 可変長数値の識別子を定義
const (
	varUint8Max        = 0xfc
	varUint16Tag uint8 = 0xfd
	varUint32Tag uint8 = 0xfe
	varUint64Tag uint8 = 0xff
)

// ServiceFlags はノードが提供するサービス一覧のビット配列
// https://en.bitcoin.it/wiki/Protocol_documentation#version
type ServiceFlags uint64

// NetPort はネットワークアドレスのポート番号を表す型
type NetPort uint16

// NetAddress は通信できるIPアドレスの型
type NetAddress struct {
	Services ServiceFlags
	IP       net.IP
	Port     NetPort
}

// DefaultUserAgent はこのクライアントのユーザーエージェント名
const DefaultUserAgent = "/btcwire:0.1.0/"
