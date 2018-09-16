package cmd

import (
	"flag"
	"fmt"
)

// Run はコマンドライン引数を受け取ってサブコマンドへディスパッチします
func Run(args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("コマンドが指定されていません")
	}
	name := args[0]
	args = args[1:]

	subCmd := dispatchCmd(name)
	if subCmd == nil {
		return fmt.Errorf("コマンドが見つかりませんでした: %s", name)
	}
	flag := flag.NewFlagSet(name, flag.ExitOnError)
	if err := subCmd.parseArgs(flag, args); err != nil {
		return err
	}

	return subCmd.exec()
}

func dispatchCmd(name string) subCmd {
	if name == "base58" {
		return &base58{}
	}
	return nil
}

// SubCmd はコマンドライン引数を受け取り実際にコマンドを実行するインターフェース
type subCmd interface {
	// 引数をパースする
	parseArgs(flag *flag.FlagSet, args []string) error
	// 実行する
	exec() error
}
