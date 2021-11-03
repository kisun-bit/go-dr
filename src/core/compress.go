package core

import (
	"fmt"
	"github.com/pierrec/lz4"
)

var ht = make([]int, 1<<16)

// Lz4Compress Lz4 Compresses data
func Lz4Compress(data []byte) (_ []byte, err error) {

	target := make([]byte, lz4.CompressBlockBound(len(data)))
	if size, err := lz4.CompressBlock(data, target, ht); err != nil || size <= 0 {
		return nil, fmt.Errorf("failed to compress: %v", err)
	} else {
		return target[:size], nil
	}
}

// Lz4Decompress Lz4 Decompress data
func Lz4Decompress(data []byte, len_ int64) (_ []byte, err error) {

	tmp := make([]byte, len_)
	if n, err := lz4.UncompressBlock(data, tmp); err != nil {
		return nil, err
	} else {
		return tmp[:n], nil
	}
}
