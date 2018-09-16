package core

// NetworkType はBTCのネットワークタイプを表す型
type NetworkType byte

const (
	// MainNetwork はメインネットワークを表す値
	MainNetwork NetworkType = 0x00
	// TestNetwork はテストネットワークを表す値
	TestNetwork NetworkType = 0x6f
)
