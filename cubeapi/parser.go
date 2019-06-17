package cubeapi

import (
	"bytes"
	"encoding/binary"

	"github.com/pkg/errors"
)

// RespBuffer struct for parsing data
type RespBuffer struct {
	buffer *bytes.Buffer
}

// CreateRespBuffer creates RespBuffer
func CreateRespBuffer(buf []byte) *RespBuffer {
	return &RespBuffer{
		buffer: bytes.NewBuffer(buf),
	}
}

// Len returns current length of buffer
func (buf *RespBuffer) Len() int {
	return buf.buffer.Len()
}

// ParseHeader parses header
func (buf *RespBuffer) ParseHeader(h *Header) error {
	err := buf.ParseInt32(&h.SvcID)
	if err != nil {
		return errors.Wrap(err, "failed to parse header")
	}
	err = buf.ParseInt32(&h.BodyLength)
	if err != nil {
		return errors.Wrap(err, "failed to parse header")
	}
	err = buf.ParseInt32(&h.RequestID)
	if err != nil {
		return errors.Wrap(err, "failed to parse header")
	}
	return nil
}

// ParseInt32 parses int32
func (buf *RespBuffer) ParseInt32(i *int32) error {
	if buf.Len() < int32Len {
		return errors.Wrap(ErrNotEnoughData, "failed to parse int32")
	}
	*i = int32(binary.LittleEndian.Uint32(buf.buffer.Next(int32Len)))
	return nil
}

// ParseInt64 parses int64
func (buf *RespBuffer) ParseInt64(i *int64) error {
	if buf.Len() < int64Len {
		return errors.Wrap(ErrNotEnoughData, "failed to parse int64")
	}
	*i = int64(binary.LittleEndian.Uint64(buf.buffer.Next(int64Len)))
	return nil
}

// ParseString parses string
func (buf *RespBuffer) ParseString(s *string) error {
	var strLen int32
	err := buf.parseStrLen(&strLen)
	if err != nil {
		return errors.Wrap(err, "failed to parse string")
	}
	err = buf.parseStr(s, strLen)
	if err != nil {
		return errors.Wrap(err, "failed to parse string")
	}
	return nil
}

func (buf *RespBuffer) parseStrLen(strLen *int32) error {
	err := buf.ParseInt32(strLen)
	if err != nil {
		return errors.Wrap(err, "failed to parse str len")
	}
	if *strLen < 0 {
		return errors.Wrap(ErrIncorrectData, "failed to parse str len: value < 0")
	}
	return nil
}

func (buf *RespBuffer) parseStr(s *string, strLen int32) error {
	if buf.Len() < int(strLen) {
		return errors.Wrap(ErrNotEnoughData, "failed to parse str")
	}
	*s = string(buf.buffer.Next(int(strLen)))
	return nil
}
