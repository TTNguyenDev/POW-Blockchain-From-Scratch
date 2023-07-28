package main

import (
	"github.com/libp2p/go-libp2p"

	"blockchain_from_scratch/cli"
)

func main() {
	cli := cli.CLI{}
	cli.Run()
}
