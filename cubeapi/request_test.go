package cubeapi_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/Apakhov/cube/cubeapi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestWriteHeader(t *testing.T) {
	buf := cubeapi.CreateSendBuffer()
	buf.WriteHeader(0x1, 0x2)
	if !bytes.Equal(buf.Bytes(), []byte{1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0}) {
		require.Equal(t, buf.Bytes(), []byte{1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0})
	}
}

func TestWriteInt32(t *testing.T) {
	buf := cubeapi.CreateSendBuffer()
	buf.WriteInt32(0x42)
	buf.WriteHeader(0x1, 0x2)
	if !bytes.Equal(buf.Bytes(), []byte{1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0x42, 0, 0, 0}) {
		require.Equal(t, buf.Bytes(), []byte{1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0x42, 0, 0, 0}, "case 1")
	}

	buf = cubeapi.CreateSendBuffer()
	buf.WriteInt32(0x42056784)
	buf.WriteHeader(0x1, 0x2)
	if !bytes.Equal(buf.Bytes(), []byte{1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0x84, 0x67, 0x05, 0x42}) {
		require.Equal(t, buf.Bytes(), []byte{1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0x84, 0x67, 0x05, 0x42}, "case 2")
	}
}

func TestWriteString(t *testing.T) {
	buf := cubeapi.CreateSendBuffer()
	buf.WriteHeader(0x1, 0x2)
	str := string([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	err := buf.WriteString(str)
	if err != nil {
		require.Equal(t, nil, err.Error(), "expected no error")
		return
	}
	if !bytes.Equal(buf.Bytes(), []byte{1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0x8, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8}) {
		require.Equal(t, buf.Bytes(), []byte{1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0x8, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8})
	}
}

func TestWriteStringErr(t *testing.T) {
	buf := cubeapi.CreateSendBuffer()
	buf.WriteHeader(0x1, 0x2)
	err := buf.WriteString(string(make([]byte, int64(math.MaxInt32)+1, int64(math.MaxInt32)+1)))
	if err == nil || errors.Cause(err) != cubeapi.ErrStringTooLong {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		require.Equal(t, cubeapi.ErrStringTooLong.Error(), errStr, "expected error")
		return
	}
}

func TestWriteInt32OnPosErr(t *testing.T) {
	buf := cubeapi.CreateSendBuffer()
	buf.WriteHeader(0x1, 0x2)
	err := buf.WriteInt32OnPos(42, buf.Len()+100)
	if err == nil || errors.Cause(err) != cubeapi.ErrBadWritingPos {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		require.Equal(t, cubeapi.ErrBadWritingPos.Error(), errStr, "expected error")
		return
	}
}
