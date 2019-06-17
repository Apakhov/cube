package cubeapi

import (
	"encoding/binary"
	"math"

	"github.com/pkg/errors"
)

type SendBuffer struct {
	buffer []byte
}

func (buf *SendBuffer) Bytes() []byte {
	return buf.buffer
}

func (buf *SendBuffer) Len() int {
	return len(buf.buffer)
}

func (buf *SendBuffer) WriteHeader(svcID int32, bodyLen int32) {
	buf.WriteInt32OnPos(svcID, 0)
	buf.WriteInt32OnPos(bodyLen, 4)
	buf.WriteInt32OnPos(0x00000000, 8)
}

func (buf *SendBuffer) WriteInt32OnPos(i int32, pos int) error {
	if pos < 0 || buf.Len() < pos+4 {
		return errors.Wrap(ErrBadWritingPos, "can't write int32")
	}
	binary.LittleEndian.PutUint32(buf.buffer[pos:pos+4], uint32(i))
	return nil
}

func (buf *SendBuffer) WriteInt32(i int32) {
	l := buf.Len()
	buf.buffer = append(buf.buffer, 0, 0, 0, 0)
	binary.LittleEndian.PutUint32(buf.buffer[l:], uint32(i))
}

func (buf *SendBuffer) WriteString(s string) error {
	if len(s) > math.MaxInt32 {
		return errors.Wrap(ErrStringTooLong, "can't write string")
	}
	buf.WriteStrLen(int32(len(s)))
	buf.WriteStr(s)
	return nil
}

func (buf *SendBuffer) WriteStrLen(i int32) {
	buf.WriteInt32(i)
}

func (buf *SendBuffer) WriteStr(s string) {
	buf.buffer = append(buf.buffer, s...)
}

func CreateSendBuffer() *SendBuffer {
	return &SendBuffer{make([]byte, headerLen, headerLen)}
}
