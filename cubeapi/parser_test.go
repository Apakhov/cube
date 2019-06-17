package cubeapi_test

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"

	"github.com/Apakhov/cube/cubeapi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// func TestCreateRespBuffer(t *testing.T) {
// 	buf := cubeapi.CreateRespBuffer([]byte{1, 1, 1, 1})
// 	if buf.Len() != 4 {
// 		require.Equal(t, 4, buf.Len())
// 	}
// }

func TestParseHeader(t *testing.T) {
	testBytes := []byte{0x4, 0, 0, 0, 0x16, 0, 0, 0, 0x3, 0, 0, 0}

	exp := &cubeapi.Header{
		SvcID:      0x4,
		BodyLength: 0x16,
		RequestID:  0x3,
	}
	res := &cubeapi.Header{}
	buf := cubeapi.CreateRespBuffer(testBytes)
	buf.ParseHeader(res)

	err := buf.Error()
	if err != nil {
		require.Equal(t, nil, err.Error(), "expected no error")
		return
	}
	if err == nil && !reflect.DeepEqual(exp, res) {
		require.Equal(t, exp, res, "result difference")
	}
}

func TestParseHeaderSlow(t *testing.T) {
	testBytes := []byte{0x4, 0, 0, 0, 0x16, 0, 0, 0, 0x3, 0, 0, 0}

	exp := &cubeapi.Header{
		SvcID:      0x4,
		BodyLength: 0x16,
		RequestID:  0x3,
	}
	res := &cubeapi.Header{}
	buf := cubeapi.CreateRespBuffer(testBytes[:6])
	go func() {
		buf.ParseHeader(res)
		buf.End()
	}()
	buf.Write(testBytes[6:])
	buf.Finished()

	buf.Wait()

	err := buf.Error()
	if err != nil {
		require.Equal(t, nil, err.Error(), "expected no error")
		return
	}
	if err == nil && !reflect.DeepEqual(exp, res) {
		require.Equal(t, exp, res, "result difference")
	}
}

func TestParseHeaderErr(t *testing.T) {
	testBytess := [][]byte{
		{0x4, 0, 0, 0, 0x16, 0, 0, 0, 0x3, 0, 0},
		{0x4, 0, 0, 0, 0x16, 0, 0, 0},
		{0x4, 0, 0},
	}
	for i, testBytes := range testBytess {
		res := &cubeapi.Header{}
		buf := cubeapi.CreateRespBuffer(testBytes)
		go func() {
			buf.ParseHeader(res)
			buf.End()
		}()
		buf.Finished()
		buf.Wait()
		err := buf.Error()
		if err == nil || errors.Cause(err) != cubeapi.ErrNotEnoughData {
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			require.Equal(t, cubeapi.ErrNotEnoughData.Error(), errStr, fmt.Sprintf("%d expected error", i))
		}
	}
}

func buildString(str string) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(str)))
	return append(buf, str...)
}

func TestParseString(t *testing.T) {
	testStrs := []string{
		"",
		"string",
		`long string with whitespaces, 		tabs,

		new lines, ` + " and weird symbols \n\t\r\\" + "and hmuric :(",
	}
	for i, testStr := range testStrs {
		res := testStr
		buf := cubeapi.CreateRespBuffer(buildString(testStr))
		go func() {
			buf.ParseString(&res)
			buf.End()
		}()
		buf.Finished()
		buf.Wait()

		err := buf.Error()
		if err != nil {
			require.Equal(t, nil, err.Error(), "expected no error")
			return
		}
		if testStr != res {
			require.Equal(t, testStr, res, fmt.Sprintf("%d result difference", i))
		}
	}
}

func TestParseStringErr(t *testing.T) {
	testCases := []struct {
		str []byte
		err error
	}{
		{[]byte{1, 1}, cubeapi.ErrNotEnoughData},
		{[]byte{0x4, 0, 0, 0}, cubeapi.ErrNotEnoughData},
		{[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, cubeapi.ErrNotEnoughData},
		{[]byte{0x00, 0x00, 0x00, 0xF0}, cubeapi.ErrIncorrectData},
	}
	for i, testCase := range testCases {
		var res string
		buf := cubeapi.CreateRespBuffer(testCase.str)
		buf.Finished()

		go func() {
			buf.ParseString(&res)
			buf.End()
		}()
		buf.Wait()

		err := buf.Error()
		if err == nil || errors.Cause(err) != testCase.err {
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			require.Equal(t, testCase.err.Error(), errStr, fmt.Sprintf("%d expected error", i))
		}
	}
}

func TestParseInt64(t *testing.T) {
	testBytes := []byte{42, 0, 0, 0, 0, 0, 0, 0}

	exp := int64(42)
	var res int64
	buf := cubeapi.CreateRespBuffer(testBytes)
	buf.ParseInt64(&res)

	err := buf.Error()
	if err != nil {
		require.Equal(t, nil, err.Error(), "expected no error")
		return
	}
	if err == nil && !reflect.DeepEqual(exp, res) {
		require.Equal(t, exp, res, "result difference")
	}
}

func TestParseInt64Err(t *testing.T) {
	testBytes := []byte{42, 0, 0, 0, 0, 0, 0}

	exp := cubeapi.ErrNotEnoughData
	var res int64
	buf := cubeapi.CreateRespBuffer(testBytes)

	go func() {
		buf.ParseInt64(&res)
		buf.End()
	}()
	buf.Finished()
	buf.Wait()

	err := buf.Error()
	if err == nil || errors.Cause(err) != exp {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		require.Equal(t, exp.Error(), errStr, "expected error")
	}
}
