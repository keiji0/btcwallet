package main

import (
	"fmt"
	"os"

	"github.com/keiji0/btc-wallet/key"

	"github.com/keiji0/btc-wallet/cmd"
)

func main() {
	key, err := key.GenerateKey()
	if err != nil {
		panic("error")
	}

	fmt.Println(key.PublicKey().Bytes())
	os.Exit(0)

	if err := cmd.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
