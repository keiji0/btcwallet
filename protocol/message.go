package protocol

import (
	"bytes"
	"io"
	"reflect"

	"github.com/keiji0/btcwallet/core"
	"github.com/keiji0/btcwallet/util/hash"
	"github.com/pkg/errors"
)

// Message はBitcoinノードへ送信するコマンドのインターフェースになります
// 各コマンドはこのインターフェースを実装します
type Message interface {
	// Messageのコマンド名を返す
	Command() string
	// MessageのPayloadをシリアライズする
	Serialize(w io.Writer) error
	// MessageのPayloadをデシリアライズする
	Deserialize(r io.Reader) error
}

// メッセージのデータ構造
// https://en.bitcoin.it/wiki/Protocol_documentation#Message_structure

// messageCommandSize はメッセージコマンドのバイトサイズ
// 余ったバイトは0で埋められます
const messageCommandSize = 12

// checksumSize はメッセージ内のチェックサムのバイトサイズ
const messageChecksumSize = 4

// messageMaxSize はメッセージの最大サイズ
// https://github.com/bitcoin/bitcoin/blob/0.17/src/serialize.h#L27
const messageMaxSize = 0x02000000

// ビットコインノードへ送信するメッセージのヘッダー
type messageHeader struct {
	magic    MessageMagic
	command  [messageCommandSize]byte
	length   uint32
	checksum [messageChecksumSize]byte
}

// メッセージコマンドの一覧
var messages = []Message{
	&MsgVersion{},
}

// コマンド名とMessageTypeのマップ、ちょっとでも早くアクセスするため
var messageMap = map[string]reflect.Type{}

func init() {
	for _, msg := range messages {
		messageMap[msg.Command()] = reflect.TypeOf(msg).Elem()
	}
}

// Send はメッセージをネットワークに送信します
func Send(w io.Writer, netType core.NetworkType, msg Message) (err error) {
	h := messageHeader{}

	// NetTypeからmagicを取得
	h.magic, err = NetworkTypeMessageMagic(netType)
	if err != nil {
		return err
	}

	// ProtocolにあったCommandに変換
	if messageCommandSize < len(msg.Command()) {
		return errors.Errorf("command name is too long: %s", msg.Command())
	}
	copy(h.command[:], msg.Command())

	// コマンド本体のバイト列を取得
	payload := &bytes.Buffer{}
	if err := msg.Serialize(payload); err != nil {
		return err
	}

	h.length = uint32(len(payload.Bytes()))
	if messageMaxSize < h.length {
		return errors.Errorf("MessageのPayloadのサイズが規定値より大きいです: command=%s, size=%d", msg.Command(), h.length)
	}

	copy(h.checksum[:], hash.Sha256x2(payload.Bytes())[0:messageChecksumSize])

	// メッセージヘッダーを送信
	if err := BulkSerialize(w, h.magic, h.command, h.length, h.checksum); err != nil {
		return err
	}
	// Payloadを送信
	if _, err := w.Write(payload.Bytes()); err != nil {
		return errors.Wrapf(err, "Payloadの送信に失敗しました")
	}

	return nil
}

// Receive はネットワークからメッセージを受信します
func Receive(r io.Reader) (Message, error) {
	h, err := readMessageHeader(r)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, h.length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, errors.Wrapf(err, "Payloadの読み込みに失敗しました: command=%s", h.command)
	}

	checksum := hash.Sha256x2(payload)
	if !bytes.Equal(checksum, h.checksum[:]) {
		return nil, errors.Errorf("Payloadのチェックサムが一致しません: command=%s", h.command)

	}

	msg, err := newMessage(string(h.command[:]))
	if err != nil {
		return nil, err
	}

	if err := msg.Deserialize(r); err != nil {
		return nil, err
	}

	return msg, nil
}

// newMessage は指定したコマンド名のメッセージを生成します
func newMessage(command string) (Message, error) {
	t, ok := messageMap[command]
	if !ok {
		return nil, errors.Errorf("メッセージコマンドが見つかりませんでした: %s", command)
	}
	i, ok := reflect.New(t).Interface().(Message)
	if !ok {
		return nil, errors.Errorf("Message Interfaceにキャストできませんでした: %q", t)
	}
	return i, nil
}

// readMessageHeader はMessageHeaderをネットワークから読み込みます
func readMessageHeader(r io.Reader) (*messageHeader, error) {
	h := &messageHeader{}

	if err := BulkDeserialize(r, &h.magic, &h.command, &h.length, &h.checksum); err != nil {
		return nil, err
	}

	if messageMaxSize < h.length {
		return nil, errors.Errorf("MessageのPayloadのサイズが規定値より大きいです: command=%s, size=%d", h.command, h.length)
	}

	return h, nil
}
