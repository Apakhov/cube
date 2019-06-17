package oauth2

import (
	"github.com/Apakhov/cube/cubeapi"
	"github.com/pkg/errors"
)

// RespBuffer struct for parsing oauth2 response
type RespBuffer struct {
	buffer *cubeapi.RespBuffer
}

// CreateRespBuffer creates RespBuffer
func CreateRespBuffer(buf []byte) *RespBuffer {
	return &RespBuffer{
		buffer: cubeapi.CreateRespBuffer(buf),
	}
}

// Len returns current len of buffer
func (buf *RespBuffer) Len() int {
	return buf.buffer.Len()
}

// ParseOAUTH2Resp parses oauth2 response
func (buf *RespBuffer) ParseOAUTH2Resp() (r ResponseOAUTH2, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(switchError(err), "failed to parse OAUTH2 response")
		}
	}()
	h := &cubeapi.Header{}
	err = buf.buffer.ParseHeader(h)
	if err != nil {
		return
	}
	if h.SvcID != cubeOAUTH2SvcID {
		err = errors.Wrap(ErrIncorrectSVCID, "failed to parse OAUTH2 response")
		return
	}
	if int(h.BodyLength) != buf.Len() {
		err = errors.Wrap(ErrIncorrectBodyLen, "failed to parse OAUTH2 response")
		return
	}
	r, err = buf.parseOAUTH2RespBody()
	return
}

func (buf *RespBuffer) parseOAUTH2RespBody() (r ResponseOAUTH2, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(switchError(err), "failed to parse OAUTH2 response body")
		}
	}()
	err = buf.parseOAUTH2RespReturnCode(&r)
	if err != nil {
		return
	}

	if r.ReturnCode != CubeOAUTH2ErrCodeOK {
		err = buf.parseOAUTH2ErrString(&r)
		return
	}

	err = buf.parseOAUTH2ClientID(&r)
	if err != nil {
		return
	}
	err = buf.parseOAUTH2ClientType(&r)
	if err != nil {
		return
	}
	err = buf.parseOAUTH2Username(&r)
	if err != nil {
		return
	}
	err = buf.parseOAUTH2ExpiresInInfo(&r)
	if err != nil {
		return
	}
	err = buf.parseOAUTH2UserID(&r)
	return
}

func (buf *RespBuffer) parseOAUTH2RespReturnCode(r *ResponseOAUTH2) error {
	err := buf.buffer.ParseInt32(&r.ReturnCode)
	if err != nil {
		return errors.Wrap(switchError(err), "failed to parse return code")
	}
	return nil
}

func (buf *RespBuffer) parseOAUTH2ErrString(r *ResponseOAUTH2) error {
	err := buf.buffer.ParseString(&r.ErrorString)
	if err != nil {
		return errors.Wrap(switchError(err), "failed to parse error string")
	}
	return nil
}

func (buf *RespBuffer) parseOAUTH2ClientID(r *ResponseOAUTH2) error {
	err := buf.buffer.ParseString(&r.CliendID)
	if err != nil {
		return errors.Wrap(switchError(err), "failed to parse client id")
	}
	return nil
}

func (buf *RespBuffer) parseOAUTH2ClientType(r *ResponseOAUTH2) error {
	err := buf.buffer.ParseInt32(&r.ClientType)
	if err != nil {
		return errors.Wrap(switchError(err), "failed to parse client type")
	}
	return nil
}

func (buf *RespBuffer) parseOAUTH2Username(r *ResponseOAUTH2) error {
	err := buf.buffer.ParseString(&r.Username)
	if err != nil {
		return errors.Wrap(switchError(err), "failed to parse username")
	}
	return nil
}

func (buf *RespBuffer) parseOAUTH2ExpiresInInfo(r *ResponseOAUTH2) error {
	err := buf.buffer.ParseInt32(&r.ExpiresIn)
	if err != nil {
		return errors.Wrap(switchError(err), "failed to parse expires_in data")
	}
	return nil
}

func (buf *RespBuffer) parseOAUTH2UserID(r *ResponseOAUTH2) error {
	err := buf.buffer.ParseInt64(&r.UserID)
	if err != nil {
		return errors.Wrap(switchError(err), "failed to parse user id")
	}
	return nil
}
