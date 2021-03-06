package main

import (
	"fmt"
	"github.com/kisun-bit/go-dr/src/core"
)

func testCompress() {
	data := "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
		"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
		"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
		"111111111112311111111111111166661111111111111111111111111111111111444111111111111111111117511111111"
	fmt.Println(fmt.Sprintf("%10s: %v\n%10s: %v\nLen %d",
		"before", []byte(data), "", data), len([]byte(data)))
	compressedRet, _ := core.Lz4Compress([]byte(data))
	fmt.Println(fmt.Sprintf("%10s: %v\n%10s: %v\nLen %d",
		"compressed", compressedRet, "", string(compressedRet), len(compressedRet)))
	compressBefore, _ := core.Lz4Decompress(compressedRet, 1024)
	fmt.Println(fmt.Sprintf("%10s: %v\n%10s: %v\nLen %d",
		"decompress", compressBefore, "", string(compressBefore), len(compressBefore)))
}
