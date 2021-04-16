package core

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	_ "github.com/OneOfOne/xxhash"
	"github.com/cespare/xxhash"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// ############## 基于AES（CBC模式）实现加密/解密

var iv = []byte{0x31, 0x37, 0x36, 0x35, 0x34, 0x33, 0x32, 0x31, 0x38, 0x27, 0x36, 0x35, 0x33, 0x23, 0x32, 0x33}

func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AESEncrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, iv)
	data = pKCS7Padding(data, blockSize)

	var cryptCode = make([]byte, len(data))
	blockMode.CryptBlocks(cryptCode, data)
	return cryptCode, nil
}

func AESEncryptToHexString(data, key []byte) (string, error) {
	r, err := AESEncrypt(data, key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(r), nil
}

func AESDecrypt(cryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cryptData))
	blockMode.CryptBlocks(origData, cryptData)
	origData = pKCS7UnPadding(origData)
	return origData, nil
}

func AESDecryptHexStringToOrigin(hexStr string, key []byte) (string, error) {
	in, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	origin, err := AESDecrypt(in, key)
	if err != nil {
		return "", err
	}
	return string(origin), nil
}

// XXHashFile 计算任意大小文件内容(路径为path)的hashcode
//
// 通过指定minCalcBlk来确定，每一个分块的计算样本长度，例如：大文件被分成1000个分块，仅取每一个分块的前minCalcBl参与哈希计算
// 文件大小与切分分块数的对应关系如下：
// ----------------------------------
// |   file size     |   block num  |
// |   0 - 1MB       |   1          |
// |   1MB - 100MB   |   100        |
// |   > 100 MB      |   1000       |
// ----------------------------------
// 计算逻辑， 例如某文件大小为100MB(即为104857600B)，则分为1000个分块，分块的size分布为104857，104857，104857...600。
// 计算每一个分块的前minCalcBlk Byte的hashCode（注意：分块长度不足minCalcBlk，则计算整个分块），
// 得到XXHash_1, XXHash_2, XXHash_1...XXHash_1000，
// 拼接得到 `XXHash_1|XXHash_2|XXHash_1|...|XXHash_1000`，并对其做整体哈希计算得到该文件的hashcode.
func XXHashFile(path string, minCalcBlk uint64) (uint64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	fp, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer fp.Close()

	var (
		perBlkSize int64
		loopNum    int
	)

	_1mb := int64(math.Pow(1024, 2))
	fSize := fi.Size()
	if 0 < fSize && fSize <= _1mb { // 0 - 1MB
		loopNum, perBlkSize = 1, _1mb
	} else if _1mb < fSize && fSize < 100*_1mb { // 1MB - 100MB
		loopNum, perBlkSize = 100, fSize/(100-1)
	} else { // > 100MB
		loopNum, perBlkSize = 1000, fSize/(1000-1)
	}

	var count int64 = 0
	bf := make([]byte, minCalcBlk)
	hashCodes := make([]string, 0)

	for {
		realLen, err := fp.Read(bf)
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}

		count += perBlkSize
		if count <= perBlkSize*int64(loopNum-1) {
			_, _ = fp.Seek(count, 0)
		}
		hashCodes = append(hashCodes, strconv.FormatUint(xxhash.Sum64(bf[:realLen]), 10))
	}

	if len(hashCodes) == 0 {
		fileHashCode, err := strconv.ParseUint(hashCodes[0], 10, 64)
		if err != nil {
			return 0, err
		}
		return fileHashCode, nil
	} else {
		return xxhash.Sum64([]byte(strings.Join(hashCodes, "|"))), nil
	}
}
