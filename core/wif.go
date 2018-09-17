package core

// WIF Wallet Import Formatの略で、ウォレットの秘密鍵を形式の一つ
type WIF struct {
	pk *PrivateKey
}

// NewWIF はWIFを秘密鍵から生成する
func NewWIF(pk *PrivateKey) *WIF {
	return &WIF{pk: pk}
}

// String WIFの文字列として取得する
func (wif *WIF) String() string {
	return NewBase58Check(WIFVersionPrefix.Byte(), wif.pk.Bytes()).String()
}
