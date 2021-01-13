package stario

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

type StarBuffer struct {
	io.Reader
	io.Writer
	io.Closer
	datas   []byte
	pStart  int
	pEnd    int
	cap     int
	isClose bool
	isEnd   bool
	rmu     sync.Mutex
	wmu     sync.Mutex
}

func NewStarBuffer(cap int) *StarBuffer {
	rtnBuffer := new(StarBuffer)
	rtnBuffer.cap = cap
	rtnBuffer.datas = make([]byte, cap)
	return rtnBuffer
}

func (star *StarBuffer) Free() int {
	return star.cap - star.Len()
}

func (star *StarBuffer) Cap() int {
	return star.cap
}

func (star *StarBuffer) Len() int {
	length := star.pEnd - star.pStart
	if length < 0 {
		return star.cap + length - 1
	}
	return length
}

func (star *StarBuffer) getByte() (byte, error) {
	if star.isClose || (star.isEnd && star.Len() == 0) {
		return 0, io.EOF
	}
	if star.Len() == 0 {
		return 0, errors.New("no byte available now")
	}
	data := star.datas[star.pStart]
	star.pStart++
	if star.pStart == star.cap {
		star.pStart = 0
	}
	return data, nil
}

func (star *StarBuffer) putByte(data byte) error {
	if star.isClose || star.isEnd {
		return io.EOF
	}
	kariEnd := star.pEnd + 1
	if kariEnd == star.cap {
		kariEnd = 0
	}
	if kariEnd == star.pStart {
		for {
			time.Sleep(time.Microsecond)
			if kariEnd != star.pStart {
				break
			}
		}
	}
	star.datas[star.pEnd] = data
	star.pEnd = kariEnd
	return nil
}
func (star *StarBuffer) Close() error {
	star.isClose = true
	return nil
}
func (star *StarBuffer) Read(buf []byte) (int, error) {
	if star.isClose {
		return 0, io.EOF
	}
	if buf == nil {
		return 0, errors.New("buffer is nil")
	}
	star.rmu.Lock()
	defer star.rmu.Unlock()
	var sum int = 0
	for i := 0; i < len(buf); i++ {
		data, err := star.getByte()
		if err != nil {
			if err == io.EOF {
				return sum, err
			}
			return sum, nil
		}
		buf[i] = data
		sum++
	}
	return sum, nil
}

func (star *StarBuffer) Write(bts []byte) (int, error) {
	if bts == nil || star.isClose {
		star.isEnd = true
		return 0, io.EOF
	}
	star.wmu.Lock()
	defer star.wmu.Unlock()
	var sum = 0
	for i := 0; i < len(bts); i++ {
		err := star.putByte(bts[i])
		if err != nil {
			fmt.Println("Write bts err:", err)
			return sum, err
		}
		sum++
	}
	return sum, nil
}
