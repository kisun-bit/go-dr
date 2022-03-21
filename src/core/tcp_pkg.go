package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/kisun-bit/go-dr/src/log"
	"io"
	"net"
	"sync"
)

// Encode When sending a native socket packet, the packet length is added to the packet header
func EncodeSocketPkg(raw []byte) (_ []byte, err error) {
	var length = int64(len(raw))
	var __pkg = new(bytes.Buffer)

	if err = binary.Write(__pkg, binary.LittleEndian, length); err != nil {
		log.DLogger.Fmt.Errorf("failed to write message length: %v", err)
		return nil, err
	}

	__pkg.Write(raw)
	return __pkg.Bytes(), nil
}

// DecodeAndHandleSocketConn Handle native socket connections and parse a complete packet to prevent sticky packets
func DecodeAndHandleSocketConn(
	conn net.Conn,
	handler func(pkg []byte) (err error),
	initFunc func(args ...interface{}) (err error)) (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic: %v", e)
		}
		if err != nil {
			log.DLogger.Error(fmt.Sprintf("err: %v", err)) // output the trace stack
		}
	}()

	var (
		once         sync.Once
		n            int
		lengthHeader = make([]byte, 4)
		lengthPkg    int32
	)

	for {
		n, err = conn.Read(lengthHeader)
		if err == io.EOF {
			if n == 0 {
				log.DLogger.Fmt.Debugf("finished")
				break
			}
			continue
		}

		if err != nil {
			return err
		}

		if n != len(lengthHeader) {
			return errors.New("invalid pkg length")
		}

		tmpLenBuff := bytes.NewBuffer(lengthHeader)
		if err = binary.Read(tmpLenBuff, binary.LittleEndian, &lengthPkg); err != nil {
			log.DLogger.Fmt.Errorf("failed to read from length buff: %v", err)
			return err
		}

		var package_ bytes.Buffer
		var leftBytes = lengthPkg

		for {
			if leftBytes == 0 {
				once.Do(func() {
					if err = initFunc(); err != nil {
						panic(err)
					}
				})
				if err = handler(package_.Bytes()); err != nil {
					return err
				}
				break // next package
			}

			// There's not enough length
			__pkgBuff := make([]byte, leftBytes)
			__readLen, err := conn.Read(__pkgBuff)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			leftBytes -= int32(__readLen)
			package_.Write(__pkgBuff[:__readLen])
		}
	}

	return nil
}
