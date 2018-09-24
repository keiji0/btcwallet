package protocol

import (
	"testing"
)

func TestMessage(t *testing.T) {

	msgVersion, err := CreateMessage("version")
	if err != nil {
		t.Error(err)
	}
	if msgVersion.Command() != "version" {
		t.Errorf("生成したメッセージのコマンド名が一致しません")
	}
}
