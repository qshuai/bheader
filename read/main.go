package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/bcext/gcash/blockchain"
	"github.com/bcext/gcash/chaincfg"
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

	fmt.Printf("block hash: %s\n", blockheader.BlockHash())
	fmt.Printf("block version: %#x\n", blockheader.Version)
	fmt.Printf("block previous block: %s\n", blockheader.PrevBlock)
	fmt.Printf("block merkle root: %s\n", blockheader.MerkleRoot)
	fmt.Printf("block timestamp: %d | time: %s\n", blockheader.Timestamp.Unix(), blockheader.Timestamp)
	fmt.Printf("block bits: %#x | big num: %s\n", blockheader.Bits, blockchain.CompactToBig(blockheader.Bits))
	difficulty, _ := getDifficultyRatio(blockheader.Bits, &chaincfg.TestNet3Params)
	fmt.Printf("difficulty: %f\n", difficulty)
	fmt.Printf("block nonce: %#x", blockheader.Nonce)
}

// getDifficultyRatio returns the proof-of-work difficulty as a multiple of the
// minimum difficulty using the passed bits field from the header of a block.
func getDifficultyRatio(bits uint32, params *chaincfg.Params) (float64, error) {
	// The minimum difficulty is the max possible proof-of-work limit bits
	// converted back to a number.  Note this is not the same as the proof of
	// work limit directly because the block difficulty is encoded in a block
	// with the compact form which loses precision.
	max := blockchain.CompactToBig(params.PowLimitBits)
	target := blockchain.CompactToBig(bits)

	difficulty := new(big.Rat).SetFrac(max, target)
	outString := difficulty.FloatString(8)
	diff, err := strconv.ParseFloat(outString, 64)
	if err != nil {

		return 0, fmt.Errorf("cannot get difficulty: %v", err)
	}
	return diff, nil
}
