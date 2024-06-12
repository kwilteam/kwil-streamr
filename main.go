package main

import (
	"fmt"
	"os"

	"github.com/kwilteam/kwil-db/cmd/kwild/root"
	"github.com/kwilteam/kwil-streamr/extensions"
)

func init() {
	err := extensions.RegisterExtensions()
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := root.RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
