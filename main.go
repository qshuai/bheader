package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/bcext/gcash/wire"
)

var hashZero = bytes.Repeat([]byte{0}, 32)

type indexMapping struct {
	header *wire.BlockHeader
	hash   string
}

func main() {
	dbdir := flag.String("dbdir", "/root/.bitcoin/chainstate", "utxo or blockindex database dir")
	targetFile := flag.String("file", "headers.dat", "result will save to the specified file")
	prefix := flag.Int("prefix", 98, "please input leveldb key prefix")
	flag.Parse()

	dboption := &DBOption{
		FilePath:  *dbdir,
		CacheSize: 1 << 20,
	}

	dbw, err := NewDBWrapper(dboption)
	if err != nil {
		panic(err)
	}

	iter := dbw.Iterator()
	defer iter.Close()

	var indexCount int
	indexSet := make(map[string]indexMapping)
	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		if int(iter.GetKey()[0]) == *prefix && iter.GetKeySize() == 33 {
			key := iter.GetKey()
			value := iter.GetVal()
			indexCount++

			var header wire.BlockHeader
			err = header.Deserialize(bytes.NewReader(value[(len(value) - 80):]))
			if err != nil {
				panic(err)
			}

			var hash chainhash.Hash
			hash.SetBytes(key[1:])

			indexSet[header.PrevBlock.String()] = indexMapping{
				header: &header,
				hash:   hash.String(),
			}
		}
	}

	// target file "headers.data"
	file, err := os.OpenFile(*targetFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tip := hex.EncodeToString(hashZero)
	for i := 0; i < len(indexSet); i++ {
		err := indexSet[tip].header.Serialize(file)
		if err != nil {
			panic(err)
		}

		tip = indexSet[tip].hash
	}

	err = file.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println("total block index count:", indexCount)
	fmt.Println("Done!")
}
