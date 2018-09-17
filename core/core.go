package core

// NetworkType はBTCのネットワークタイプを表す型
type NetworkType byte

const (
	// MainNetwork はメインネットワークを表す値
	MainNetwork NetworkType = iota
	// TestNetwork はテストネットワークを表す値
	TestNetwork NetworkType = iota
)

// VersionPrefix はアドレスの種類を表す識別子
// ref. https://en.bitcoin.it/wiki/List_of_address_prefixes
type VersionPrefix byte

const (
	// MainNetVersionPrefix はメインネットワークで使用できるアドレス
	MainNetVersionPrefix VersionPrefix = 0x00
	// TestNetVersionPrefix はテストネットワークで使用できるアドレス
	TestNetVersionPrefix VersionPrefix = 0x6f

	// WIFVersionPrefix は秘密鍵のWallet Import形式を識別するための識別子
	WIFVersionPrefix VersionPrefix = 0x80
)

// Byte はVersionPrefixのバイトを取得する
func (vp VersionPrefix) Byte() byte {
	return byte(vp)
}
