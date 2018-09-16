package cmd

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	b58 "github.com/keiji0/btc-wallet/util/base58"
)

type base58 struct {
	decode bool
	args   []string
}

func (b *base58) parseArgs(flag *flag.FlagSet, args []string) error {
	flag.BoolVar(&b.decode, "decode", false, "デコードする")
	flag.Parse(args)
	b.args = flag.Args()
	return nil
}

func (b *base58) exec() error {
	res, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if b.decode {
		os.Stdout.Write(b58.Decode(string(res)))
	} else {
		os.Stdout.WriteString(b58.Encode(res))
	}
	return nil
}
