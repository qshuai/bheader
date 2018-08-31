package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/bcext/gcash/wire"
)

// 读取连续的header文件
//func main() {
//	f, err := os.OpenFile("./headers", os.O_RDONLY, 0)
//	if err != nil {
//		panic(err)
//	}
//	receiver := make([]byte, 80)
//	f.Seek(160, 0)
//	n, err := f.Read(receiver)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("read %d bytes", n)
//
//	var header wire.BlockHeader
//	err = header.Deserialize(bytes.NewReader(receiver))
//	if err != nil {
//		panic(err)
//	}
//
//	spew.Dump(header)
//}

var hashZero = bytes.Repeat([]byte{0}, 32)

type indexMapping struct {
	header *wire.BlockHeader
	prev   string
}

func main() {
	dbdir := flag.String("dbdir", "/root/.bitcoin/chainstate", "utxo or blockindex database dir")
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
				prev:   hash.String(),
			}
		}
	}

	buf := bytes.NewBuffer(make([]byte, 80*1240000))
	tip := hex.EncodeToString(hashZero)
	for i := 0; i < len(indexSet); i++ {
		err := indexSet[tip].header.Serialize(buf)
		if err != nil {
			panic(err)
		}

		tip = indexSet[tip].prev
	}

	err = ioutil.WriteFile("headers", buf.Bytes(), os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println("total block index count:", indexCount)
}
