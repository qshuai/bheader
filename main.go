package main

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/bcext/gcash/wire"
	"github.com/davecgh/go-spew/spew"
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
	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		if int(iter.GetKey()[0]) == *prefix {
			value := iter.GetVal()
			fmt.Println(iter.GetKey())
			fmt.Println(iter.GetKeySize())
			fmt.Println(value)
			fmt.Println(iter.GetValSize())
			indexCount++

			var header wire.BlockHeader
			err = header.Deserialize(bytes.NewReader(value[(len(value) - 80):]))
			if err != nil {
				panic(err)
			}
			spew.Dump(header)
			if indexCount >= 3 {
				break
			}
		}
	}

	fmt.Println("total block index count:", indexCount)
}
