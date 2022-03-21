package main

import (
	"fmt"
	"github.com/kisun-bit/go_dr/src/core"
)

func testAes() {
	key := []byte("1234567890123456")
	initVal := []byte("##########!!!!!!!!!!")
	fmt.Println("before   : ", string(initVal))
	ret, _ := core.AESEncryptToHexString(initVal, key)
	fmt.Println("encrypted: ", ret)

	origin, _ := core.AESDecryptHexStringToOrigin(ret, key)
	fmt.Println("decrypted: ", origin)
}

//func main() {
//	testAes()
//}
