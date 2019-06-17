package oauth2_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/Apakhov/cube/cubeapi/oauth2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCreateOAUTH2Request(t *testing.T) {
	testBytes := []byte{2, 0, 0, 0, 22, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 5, 0, 0, 0, 116, 111, 107, 101, 110, 5, 0, 0, 0, 115, 99, 111, 112, 101}

	buf, err := oauth2.CreateOAUTH2Request("token", "scope")
	if err != nil {
		require.Equal(t, nil, err.Error(), "expected no error")
		return
	}
	if !bytes.Equal(buf.Bytes(), testBytes) {
		require.Equal(t, buf.Bytes(), testBytes)
	}
}

func TestCreateOAUTH2RequestErr(t *testing.T) {
	_, err := oauth2.CreateOAUTH2Request(string(make([]byte, int64(math.MaxInt32)+1, int64(math.MaxInt32)+1)), "test")
	if err == nil || errors.Cause(err) != oauth2.ErrStringTooLong {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		t.Logf("%+v <-> %+v", errors.Cause(err), oauth2.ErrStringTooLong)
		require.Equal(t, oauth2.ErrStringTooLong.Error(), errStr, "expected error 1")
	}
	_, err = oauth2.CreateOAUTH2Request("test", string(make([]byte, int64(math.MaxInt32)+1, int64(math.MaxInt32)+1)))
	if err == nil || errors.Cause(err) != oauth2.ErrStringTooLong {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		require.Equal(t, oauth2.ErrStringTooLong.Error(), errStr, "expected error 2")
	}
}
