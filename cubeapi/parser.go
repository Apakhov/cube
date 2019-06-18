package cubeapi

import (
	"bytes"
	"encoding/binary"

	"github.com/pkg/errors"
)

// RespBuffer struct for parsing data. It can be used in sync and async mode
//
// sync mode - straightforward using of Parse commands
//
// async mode - using Parse commands in separate go-routine
type RespBuffer struct {
	buffer         *bytes.Buffer
	bytesAvailable int64
	parseLimit     int64
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
		bytesAvailable: int64(len(buf)),
		parseLimit:     0,
	}
	return res
}

// IncreaseParseLim increases limit for parsing
func (buf *RespBuffer) IncreaseParseLim(i int64) {
	buf.parseLimit += i
}

// GetParseLim returns current parse limit
func (buf *RespBuffer) GetParseLim() int64 {
	return buf.parseLimit
}

func (buf *RespBuffer) Write(part []byte) {
	buf.finLock <- struct{}{}
	buf.buffer.Write(part)
	buf.bytesAvailable += int64(len(part))
	<-buf.finLock
}

// Finished should be called after writing all data
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

func (buf *RespBuffer) primalErrorCheck(length int64, msg string) (ok bool) {
	buf.parseLimit -= length
	if buf.parseLimit < 0 {
		buf.parseLimit = 0
		buf.createError(ErrIncorrectLen, msg)
		return false
	}
	if buf.err != nil || !buf.blockForBytes(length) {
		buf.createError(ErrNotEnoughData, msg)
		return false
	}
	return true
}

func (buf *RespBuffer) loadError(msg string) (written bool) {
	buf.err = errors.Wrap(buf.err, msg)
	return buf.err != nil
}

func (buf *RespBuffer) Error() error {
	return buf.err
}

// End should be called after all parsing commands in async mode
func (buf *RespBuffer) End() {
	buf.end <- struct{}{}
	close(buf.end)
}

// Wait
func (buf *RespBuffer) Wait() {
	<-buf.end
}

func (buf *RespBuffer) WaitChan() <-chan struct{} {
	return buf.end
}

func (buf *RespBuffer) blockForBytes(amount int64) bool {
	for {
		buf.finLock <- struct{}{}
		if int64(amount) <= buf.bytesAvailable {
			buf.bytesAvailable -= amount
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
	return
}

// ParseInt32 parses int32
func (buf *RespBuffer) ParseInt32(i *int32) {
	if buf.primalErrorCheck(int32Len, "failed to parse int32") {
		*i = int32(binary.LittleEndian.Uint32(buf.buffer.Next(int32Len)))
	}
}

// ParseInt64 parses int64
func (buf *RespBuffer) ParseInt64(i *int64) {
	if buf.primalErrorCheck(int64Len, "failed to parse int64") {
		*i = int64(binary.LittleEndian.Uint64(buf.buffer.Next(int64Len)))
	}
}

// ParseString parses string
func (buf *RespBuffer) ParseString(s *string) {
	var strLen int32
	buf.parseStrLen(&strLen)
	buf.parseStr(s, strLen)
	buf.loadError("failed to parse string")
	return
}

func (buf *RespBuffer) parseStrLen(strLen *int32) {
	buf.ParseInt32(strLen)
	if buf.loadError("failed to parse str len") {
		return
	}
	if *strLen < 0 {
		buf.createError(ErrIncorrectData, "failed to parse str len: value < 0")
	}
	return
}

func (buf *RespBuffer) parseStr(s *string, strLen int32) {
	if buf.primalErrorCheck(int64(strLen), "failed to parse str") {
		*s = string(buf.buffer.Next(int(strLen)))
	}
}
