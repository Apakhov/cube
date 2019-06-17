package cubeapi

import (
	"bytes"
	"encoding/binary"
	"sync/atomic"

	"github.com/pkg/errors"
)

// RespBuffer struct for parsing data
type RespBuffer struct {
	buffer         *bytes.Buffer
	bytesAvailable *int64
	finished       chan struct{}
	end            chan struct{}
	err            error
	finLock        chan struct{}
}

// CreateRespBuffer creates RespBuffer
func CreateRespBuffer(buf []byte) *RespBuffer {
	res := &RespBuffer{
		buffer:         bytes.NewBuffer(buf),
		finished:       make(chan struct{}, 1),
		end:            make(chan struct{}, 1),
		finLock:        make(chan struct{}, 1),
		bytesAvailable: new(int64),
	}
	atomic.StoreInt64(res.bytesAvailable, int64(len(buf)))
	return res
}

func (buf *RespBuffer) Write(part []byte) {
	buf.finLock <- struct{}{}
	buf.buffer.Write(part)
	atomic.AddInt64(buf.bytesAvailable, int64(len(part)))
	<-buf.finLock
}

func (buf *RespBuffer) Finished() {
	buf.finLock <- struct{}{}
	buf.finished <- struct{}{}
	<-buf.finLock
}

func (buf *RespBuffer) createError(err error, msg string) {
	if buf.err == nil {
		buf.err = errors.Wrap(err, msg)
	}
}

func (buf *RespBuffer) loadError(msg string) (written bool) {
	buf.err = errors.Wrap(buf.err, "msg")
	return buf.err != nil
}

func (buf *RespBuffer) Error() error {
	return buf.err
}

func (buf *RespBuffer) End() {
	buf.end <- struct{}{}
	close(buf.end)
}

func (buf *RespBuffer) Wait() {
	<-buf.end
}

func (buf *RespBuffer) WaitChan() <-chan struct{} {
	return buf.end
}

func (buf *RespBuffer) blockForBytes(amount int) bool {
	for {
		buf.finLock <- struct{}{}

		av := atomic.LoadInt64(buf.bytesAvailable)
		if int64(amount) <= av {
			atomic.AddInt64(buf.bytesAvailable, -int64(amount))
			<-buf.finLock
			return true
		}
		select {
		case <-buf.finished:
			<-buf.finLock
			return false
		default:
		}
		<-buf.finLock
	}
}

// ParseHeader parses header
func (buf *RespBuffer) ParseHeader(h *Header) {
	buf.ParseInt32(&h.SvcID)
	buf.ParseInt32(&h.BodyLength)
	buf.ParseInt32(&h.RequestID)
	buf.loadError("failed to parse header")

}

// // ParseByte parses byte
// func (buf *RespBuffer) ParseByte(b *byte) {
// 	if buf.err != nil || !buf.blockForBytes(1) {
// 		buf.createError(ErrNotEnoughData, "failed to parse byte")
// 		return
// 	}
// 	*b = buf.buffer.Next(1)[0]
// }

// ParseInt32 parses int32
func (buf *RespBuffer) ParseInt32(i *int32) {
	if buf.err != nil || !buf.blockForBytes(int32Len) {
		buf.createError(ErrNotEnoughData, "failed to parse int32")
		return
	}
	*i = int32(binary.LittleEndian.Uint32(buf.buffer.Next(int32Len)))

}

// ParseInt64 parses int64
func (buf *RespBuffer) ParseInt64(i *int64) {
	if buf.err != nil || !buf.blockForBytes(int64Len) {
		buf.createError(ErrNotEnoughData, "failed to parse int64")
		return
	}
	*i = int64(binary.LittleEndian.Uint64(buf.buffer.Next(int64Len)))

}

// ParseString parses string
func (buf *RespBuffer) ParseString(s *string) {
	var strLen int32
	buf.parseStrLen(&strLen)
	buf.parseStr(s, strLen)
	buf.loadError("failed to parse string")

}

func (buf *RespBuffer) parseStrLen(strLen *int32) {
	buf.ParseInt32(strLen)
	if buf.loadError("failed to parse str len") {
		return
	}

	if *strLen < 0 {
		buf.createError(ErrIncorrectData, "failed to parse str len: value < 0")
	}
}

func (buf *RespBuffer) parseStr(s *string, strLen int32) {
	if buf.err != nil || !buf.blockForBytes(int(strLen)) {
		buf.createError(ErrNotEnoughData, "failed to parse str")
		return
	}
	*s = string(buf.buffer.Next(int(strLen)))
}

// str := make([]byte, int(strLen))
// 	for i := 0; i < int(strLen); i++ {
// 		buf.ParseByte(&str[i])
// 		if buf.loadError("failed to parse str") {
// 			return
// 		}
// 	}

// 	*s = string(str)
