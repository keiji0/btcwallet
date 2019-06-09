package protocol

import (
	"io"
	"time"
)

// MsgVersion はノードへの接続時に通信するためのバージョンメッセージを表す型
type MsgVersion struct {
	// プロトコルのバージョン
	ProtocolVersion Version
	// ノードが提供するサービス一覧
	Services ServiceFlags
	// メッセージが作られた時刻
	Timestamp Int64Time
	// 受けてのネットワークアドレス
	AddrRerv NetAddress
	// ノードのネットワークアドレス
	AddrFrom NetAddress
	// 送信時にランダムに生成される値
	Nonce uint64
	// ユーザーエージェント名
	UserAgent UserAgentName
	// 自身のノードが持っているブロックの高さ
	StartHeight int32
	// INVを送られないようにする設定
	Relay bool
}

// Command はこのメッセージのコマンド名を返します
func (v *MsgVersion) Command() string {
	return "version"
}

// NewMsgVersion はMsgVersionメッセージを生成します
func NewMsgVersion(rervAddr, fromAddr *NetAddress, nonce uint64, blockHeight int32) *MsgVersion {
	v := &MsgVersion{
		ProtocolVersion: CurrentVersion,
		Services:        0,
		Timestamp:       Int64Time(time.Unix(time.Now().Unix(), 0)),
		AddrRerv:        *rervAddr,
		AddrFrom:        *fromAddr,
		Nonce:           nonce,
		UserAgent:       DefaultUserAgent,
		StartHeight:     blockHeight,
		Relay:           false,
	}
	return v
}

// Serialize はMessageのPayloadをシリアライズする
func (v *MsgVersion) Serialize(w io.Writer) error {
	return BulkSerialize(
		w,
		v.ProtocolVersion,
		v.Services,
		v.Timestamp,
		v.AddrRerv,
		v.AddrFrom,
		v.Nonce,
		v.UserAgent,
		v.StartHeight,
		v.Relay,
	)
}

// Deserialize はMessageのPayloadをデシリアライズする
func (v *MsgVersion) Deserialize(r io.Reader) error {
	return BulkDeserialize(
		r,
		&v.ProtocolVersion,
		&v.Services,
		&v.Timestamp,
		&v.AddrRerv,
		&v.AddrFrom,
		&v.Nonce,
		&v.UserAgent,
		&v.StartHeight,
		&v.Relay,
	)
}
