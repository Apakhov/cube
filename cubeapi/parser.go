package cubeapi

import (
	"bytes"
	"encoding/binary"

	"github.com/pkg/errors"
)

type header struct {
	svcID      int32
	bodyLength int32
	requestID  int32
}

const int32Len = 4
const int64Len = 8

func parseHeader(buffer *bytes.Buffer) (h header, err error) {
	if buffer.Len() < headerLen {
		err = errors.New("not enough data to parse header")
		return
	}
	h.svcID = int32(binary.LittleEndian.Uint32(buffer.Next(4)))
	h.bodyLength = int32(binary.LittleEndian.Uint32(buffer.Next(4)))
	h.requestID = int32(binary.LittleEndian.Uint32(buffer.Next(4)))
	return
}

func parseInt32(i *int32, buffer *bytes.Buffer) error {
	if buffer.Len() < int32Len {
		return errors.New("not enough data to parse int32")
	}
	*i = int32(binary.LittleEndian.Uint32(buffer.Next(int32Len)))
	return nil
}

func parseInt64(i *int64, buffer *bytes.Buffer) error {
	if buffer.Len() < int64Len {
		return errors.New("not enough data to parse int64")
	}
	*i = int64(binary.LittleEndian.Uint64(buffer.Next(int64Len)))
	return nil
}

func parseString(s *string, buffer *bytes.Buffer) error {
	var strLen int32
	err := parseStrLen(&strLen, buffer)
	if err != nil {
		return errors.Wrap(err, "failed to parse string")
	}
	err = parseStr(s, buffer, strLen)
	if err != nil {
		return errors.Wrap(err, "failed to parse string")
	}
	return nil
}

func parseStrLen(strLen *int32, buffer *bytes.Buffer) error {
	err := parseInt32(strLen, buffer)
	if err != nil {
		return errors.Wrap(err, "failed to parse str len")
	}
	if *strLen < 0 {
		return errors.New("failed to parse str len: value < 0")
	}
	return nil
}

func parseStr(s *string, buffer *bytes.Buffer, strLen int32) error {
	if buffer.Len() < int(strLen) {
		return errors.New("failed to parse str: not enough data")
	}
	*s = string(buffer.Next(int(strLen)))
	return nil
}
