package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/bcext/gcash/wire"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	f, err := os.OpenFile("./headers", os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	receiver := make([]byte, 80)
	f.Seek(240, 0)
	n, err := f.Read(receiver)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read %d bytes", n)

	var header wire.BlockHeader
	err = header.Deserialize(bytes.NewReader(receiver))
	if err != nil {
		panic(err)
	}

	spew.Dump(header)
}
