package core

import (
	"encoding/binary"
	"github.com/kisun-bit/go_dr/src/log"
	"io"
	"os"
	"sync"
)

type SimpleBit struct {
	Offset int64
	Length int
	Hash   uint64
}

type SimpleBitmapFile struct {
	sync.RWMutex

	path     string
	fp       *os.File
	isOpened bool
	bits     chan *SimpleBit
}

func (b *SimpleBitmapFile) GetBlockSize() int {
	return 8
}

func (b *SimpleBitmapFile) WriteByOffset(offset int64, hash uint64) (err error) {
	b.Lock()
	defer b.Unlock()

	if err = b.initHandle(); err != nil {
		return
	}

	tmpBuf := make([]byte, b.GetBlockSize())
	binary.BigEndian.PutUint64(tmpBuf, hash)

	if _, err = b.fp.WriteAt(tmpBuf, offset); err != nil {
		return
	}

	return nil
}

func (b *SimpleBitmapFile) WriteByIndex(index int64) (err error) {
	_ = index
	panic("TODO") // TODO
}

func (b *SimpleBitmapFile) ReadByOffset(offset int64) (hash uint64, err error) {
	b.Lock()
	defer b.Unlock()

	if err := b.initHandle(); err != nil {
		return 0, err
	}

	var buf = make([]byte, b.GetBlockSize())
	if _, err = b.fp.ReadAt(buf, offset); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(buf), nil
}

func (b *SimpleBitmapFile) Iter() {
	if err := b.initHandle(); err != nil {
		return
	}

	go b.__iter(0B00000000)
}

func (b *SimpleBitmapFile) __iter(__start int64) {
	b.Lock()

	defer b.Unlock()
	defer close(b.bits)

	var count int64

	buf := make([]byte, b.GetBlockSize())
	for true {
		n, err := b.fp.Read(buf)
		if err == io.EOF {
			log.DLogger.Fmt.Debugf("Bitmap read finished")
			return
		} else if err != nil {
			log.DLogger.Fmt.Errorf("Bitmap read err: %v", err)
			return
		} else {
			var _bit = new(SimpleBit)
			_bit.Offset = __start + count*int64(b.GetBlockSize())
			_bit.Length = b.GetBlockSize()
			_bit.Hash = binary.BigEndian.Uint64(buf[:n])
			b.bits <- _bit
			count++
		}
	}
}

func (b *SimpleBitmapFile) initHandle() (err error) {
	if !b.isOpened {

		b.fp, err = os.OpenFile(b.path, os.O_RDWR|os.O_CREATE, 0)
		if err != nil {

		}
		log.DLogger.Fmt.Debugf("---------init handle successfully: %v--------", b.path)
		b.isOpened = true
	}

	return
}

func (b *SimpleBitmapFile) Close() (err error) {
	if err = b.fp.Close(); err != nil {
		log.DLogger.Fmt.Warnf("close `%v` Failed: %v", b.path, err)
	}

	b.isOpened = false
	log.DLogger.Fmt.Debugf("---------close handle successfully: %v--------", b.path)
	return nil
}
