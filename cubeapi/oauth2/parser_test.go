package oauth2_test

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Apakhov/cube/cubeapi/oauth2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func buildString(str string) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(str)))
	return append(buf, str...)
}

func buildInt32(i int32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(i))
	return buf
}

func buildInt64(i int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(i))
	return buf
}

func TestParseOAUTH2RespOKCode(t *testing.T) {
	testBytes := []byte{}
	testBytes = append(testBytes, buildInt32(0x2)...)                 // svc id
	testBytes = append(testBytes, buildInt32(0x0)...)                 // body len, unknown yet
	testBytes = append(testBytes, buildInt32(0x1)...)                 // request id
	testBytes = append(testBytes, buildInt32(0x0)...)                 // return code
	testBytes = append(testBytes, buildString("test_client_id")...)   // client id
	testBytes = append(testBytes, buildInt32(2002)...)                // client type
	testBytes = append(testBytes, buildString("testuser@mail.ru")...) // username
	testBytes = append(testBytes, buildInt32(3600)...)                // expires in
	testBytes = append(testBytes, buildInt64(101010)...)              // user id
	// inserting body length
	binary.LittleEndian.PutUint32(testBytes[4:8], uint32(len(testBytes)-12))

	exp := oauth2.ResponseOAUTH2{
		ReturnCode:  oauth2.CubeOAUTH2ErrCodeOK,
		CliendID:    "test_client_id",
		ClientType:  2002,
		Username:    "testuser@mail.ru",
		ExpiresIn:   3600,
		UserID:      101010,
		ErrorString: "",
	}
	buf := oauth2.CreateRespBuffer(testBytes)
	var res oauth2.ResponseOAUTH2
	buf.ParseOAUTH2Resp(&res)

	err := buf.Error()

	if err != nil {
		require.Equal(t, nil, err.Error(), "expected no error")
		return
	}
	if err == nil && !reflect.DeepEqual(exp, res) {
		require.Equal(t, exp, res, "result difference")
	}
}

func TestParseOAUTH2RespOKCodeParts(t *testing.T) {
	wg := &sync.WaitGroup{}
	for p := 0; p < 500; p++ {
		wg.Add(1)
		go func() {
			testBytes := []byte{}

			buf := oauth2.CreateRespBuffer(testBytes)
			testBytes = append(testBytes, buildInt32(0x2)...)                 // svc id
			testBytes = append(testBytes, buildInt32(0x0)...)                 // body len, unknown yet
			testBytes = append(testBytes, buildInt32(0x1)...)                 // request id
			testBytes = append(testBytes, buildInt32(0x0)...)                 // return code
			testBytes = append(testBytes, buildString("test_client_id")...)   // client id
			testBytes = append(testBytes, buildInt32(2002)...)                // client type
			testBytes = append(testBytes, buildString("testuser@mail.ru")...) // username
			testBytes = append(testBytes, buildInt32(3600)...)                // expires in
			testBytes = append(testBytes, buildInt64(101010)...)              // user id
			// inserting body length
			binary.LittleEndian.PutUint32(testBytes[4:8], uint32(len(testBytes)-12))

			var res oauth2.ResponseOAUTH2

			exp := oauth2.ResponseOAUTH2{
				ReturnCode:  oauth2.CubeOAUTH2ErrCodeOK,
				CliendID:    "test_client_id",
				ClientType:  2002,
				Username:    "testuser@mail.ru",
				ExpiresIn:   3600,
				UserID:      101010,
				ErrorString: "",
			}

			go func() {
				buf.ParseOAUTH2Resp(&res)
				buf.End()
			}()

			for l := 0; l < len(testBytes); l += 10 {
				time.Sleep(time.Millisecond * 50)
				r := l + 10
				if r > len(testBytes) {
					r = len(testBytes)
				}
				buf.Write(testBytes[l:r])
			}
			buf.Finished()

			buf.Wait()

			err := buf.Error()

			if err != nil {
				require.Equal(t, nil, err.Error(), "expected no error")
				return
			}

			if err == nil && !reflect.DeepEqual(exp, res) {
				require.Equal(t, exp, res, "result difference")
				return
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestParseOAUTH2RespErrCode(t *testing.T) {
	wg := &sync.WaitGroup{}
	for p := 0; p < 100; p++ {
		wg.Add(1)
		go func() {
			testBytes := []byte{}

			buf := oauth2.CreateRespBuffer(testBytes)
			//header
			testBytes = append(testBytes, buildInt32(0x2)...) // svc id
			testBytes = append(testBytes, buildInt32(0x0)...) // body len, unknown yet
			testBytes = append(testBytes, buildInt32(0x1)...) // request id
			//body
			testBytes = append(testBytes, buildInt32(oauth2.CubeOAUTH2ErrCodeBadClient)...) // return code
			testBytes = append(testBytes, buildString("lol you died")...)                   // error string
			// inserting body length
			binary.LittleEndian.PutUint32(testBytes[4:8], uint32(len(testBytes)-12))
			var res oauth2.ResponseOAUTH2

			exp := oauth2.ResponseOAUTH2{
				ReturnCode: oauth2.CubeOAUTH2ErrCodeBadClient,

				ErrorString: "lol you died",
			}

			go func() {
				buf.ParseOAUTH2Resp(&res)
				buf.End()
			}()

			for l := 0; l < len(testBytes); l += 2 {
				time.Sleep(time.Millisecond * 50)
				r := l + 2
				if r > len(testBytes) {
					r = len(testBytes)
				}
				buf.Write(testBytes[l:r])
			}
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
			wg.Done()
		}()
	}
	wg.Wait()
}

func flat(bss ...[]byte) []byte {
	r := []byte{}
	for _, bs := range bss {
		r = append(r, bs...)
	}
	return r
}

type parseOAUTH2RespErrCase struct {
	bytes  []byte
	err    error
	blCorr bool
}

var parseOAUTH2RespErrCases = []parseOAUTH2RespErrCase{
	{
		bytes: flat(
			buildInt32(0x2),
			buildInt32(0x1), // incorrect body length
			buildInt32(0x1),
		),
		err: oauth2.ErrIncorrectLen,
	},
	{
		bytes: []byte{}, // no header
		err:   oauth2.ErrNotEnoughData,
	},
	{
		bytes: flat(
			buildInt32(0x43534534), // incorrect svcID
			buildInt32(0x1),
			buildInt32(0x1),
		),
		err: oauth2.ErrIncorrectSVCID,
	},
	{
		bytes: flat( // no body
			buildInt32(0x2),
			buildInt32(0x0),
			buildInt32(0x1),
		),
		err: oauth2.ErrIncorrectLen,
	},
	{
		bytes: flat(
			buildInt32(0x2),
			buildInt32(0x0),
			buildInt32(0x1),
			buildInt32(0x2), // incorrect svc error response body (no error string)
		),
		err:    oauth2.ErrIncorrectLen,
		blCorr: true,
	},
	{
		bytes: flat(
			buildInt32(0x2),
			buildInt32(0x0),
			buildInt32(0x1),
			buildInt32(0x0), // incorrect svc ok response body (no client id)
		),
		err:    oauth2.ErrIncorrectLen,
		blCorr: true,
	},
	{
		bytes: flat(
			buildInt32(0x2),
			buildInt32(0x0),
			buildInt32(0x1),
			buildInt32(0x0),
			buildString("client id"), // client id
			// incorrect svc ok response body (no client type)
		),
		err:    oauth2.ErrIncorrectLen,
		blCorr: true,
	},
	{
		bytes: flat(
			buildInt32(0x2),
			buildInt32(0x0),
			buildInt32(0x1),
			buildInt32(0x0),
			buildString("client id"), // client id
			buildInt32(0x0),          // client type
			buildInt32(0x32),         // incorrect svc ok response body (incorrect username)
		),
		err:    oauth2.ErrIncorrectLen,
		blCorr: true,
	},
	{
		bytes: flat(
			buildInt32(0x2),
			buildInt32(0x0),
			buildInt32(0x1),
			buildInt32(0x0),
			buildString("client id"), // client id
			buildInt32(0x0),          // client type
			buildString("username"),  // incorrect svc ok response body (incorrect username)
			// incorrect svc ok response body (no expires in)
		),
		err:    oauth2.ErrIncorrectLen,
		blCorr: true,
	},
	{
		bytes: flat(
			buildInt32(0x2),
			buildInt32(0x0),
			buildInt32(0x1),
			buildInt32(0x0),
			buildString("client id"), // client id
			buildInt32(0x0),          // client type
			buildString("username"),  // incorrect svc ok response body (incorrect username)
			buildInt32(0x35),         // expires in
			buildInt32(0x35),         // incorrect svc ok response body (incorrect user id)
		),
		err:    oauth2.ErrIncorrectLen,
		blCorr: true,
	},
}

func TestParseOAUTH2RespErr(t *testing.T) {
	for i, c := range parseOAUTH2RespErrCases {
		testBytes := c.bytes
		if c.blCorr {
			// inserting body length
			binary.LittleEndian.PutUint32(testBytes[4:8], uint32(len(testBytes)-12))

		}
		buf := oauth2.CreateRespBuffer(testBytes)
		buf.Finished()
		var res oauth2.ResponseOAUTH2
		buf.ParseOAUTH2Resp(&res)

		err := buf.Error()
		if err == nil || errors.Cause(err) != c.err {
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			require.Equal(t, c.err.Error(), errStr, fmt.Sprintf("%d expected error ", i))
		}
	}
}
