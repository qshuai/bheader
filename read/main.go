package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bcext/gcash/wire"
	"github.com/qshuai/tcolor"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println(tcolor.WithColor(tcolor.Red, "parameter mismatch"))
		return
	}

	file, err := os.Open(args[1])
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "open blockheader file failed"))
		return
	}

	blockheight, err := strconv.Atoi(args[2])
	if err != nil || blockheight < 0 {
		fmt.Println(tcolor.WithColor(tcolor.Red, "blockheight parameter error"))
		return
	}
	_, err = file.Seek(80*int64(blockheight), 0)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "file seek failed"))
		return
	}

	var blockheader wire.BlockHeader
	err = blockheader.Deserialize(file)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "blockheader information error"))
		return
	}

	fmt.Printf("block version: %#x\n", blockheader.Version)
	fmt.Printf("block previous block: %s\n", blockheader.PrevBlock)
	fmt.Printf("block merkle root: %s\n", blockheader.MerkleRoot)
	fmt.Printf("block time: %s\n", blockheader.Timestamp)
	fmt.Printf("block bits: %#x\n", blockheader.Bits)
	fmt.Printf("block nonce: %#x\n", blockheader.Nonce)
}
