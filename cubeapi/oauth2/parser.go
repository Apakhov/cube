package oauth2

import (
	"github.com/Apakhov/cube/cubeapi"
	"github.com/pkg/errors"
)

// RespBuffer struct for parsing oauth2 response
type RespBuffer struct {
	buffer *cubeapi.RespBuffer
	err    error
}

// CreateRespBuffer creates RespBuffer
func CreateRespBuffer(buf []byte) *RespBuffer {
	return &RespBuffer{
		buffer: cubeapi.CreateRespBuffer(buf),
	}
}

func (buf *RespBuffer) Write(part []byte) {
	buf.buffer.Write(part)
}

func (buf *RespBuffer) Finished() {
	buf.buffer.Finished()
}

func (buf *RespBuffer) createError(err error, msg string) {
	if buf.err == nil {
		buf.err = errors.Wrap(err, msg)
	}
}

func (buf *RespBuffer) loadError(msg string) (written bool) {
	if buf.buffer.Error() != nil && buf.err == nil {
		buf.err = switchError(buf.buffer.Error())
	}
	buf.err = errors.Wrap(buf.err, msg)
	return buf.err != nil
}

func (buf *RespBuffer) Error() error {
	return buf.err
}

func (buf *RespBuffer) End() {
	buf.buffer.End()
}

func (buf *RespBuffer) Wait() {
	buf.buffer.Wait()
}

func (buf *RespBuffer) WaitChan() <-chan struct{} {
	return buf.buffer.WaitChan()
}

// ParseOAUTH2Resp parses oauth2 response
func (buf *RespBuffer) ParseOAUTH2Resp(r *ResponseOAUTH2) {
	h := &cubeapi.Header{}
	buf.buffer.ParseHeader(h)
	buf.loadError("failed to parse OAUTH2 response")
	if buf.err != nil {
		return
	}
	if h.SvcID != cubeOAUTH2SvcID {
		buf.createError(ErrIncorrectSVCID, "failed to parse OAUTH2 response")
		return
	}
	buf.parseOAUTH2RespBody(r)
	buf.loadError("failed to parse OAUTH2 response")
	return
}

func (buf *RespBuffer) parseOAUTH2RespBody(r *ResponseOAUTH2) {
	if buf.err != nil {
		return
	}
	buf.parseOAUTH2RespReturnCode(r)
	if r.ReturnCode != CubeOAUTH2ErrCodeOK {
		buf.parseOAUTH2ErrString(r)
	} else {
		buf.parseOAUTH2ClientID(r)
		buf.parseOAUTH2ClientType(r)
		buf.parseOAUTH2Username(r)
		buf.parseOAUTH2ExpiresInInfo(r)
		buf.parseOAUTH2UserID(r)
	}
	buf.loadError("failed to parse OAUTH2 response body")
	return
}

func (buf *RespBuffer) parseOAUTH2RespReturnCode(r *ResponseOAUTH2) {
	if buf.err != nil {
		return
	}
	buf.buffer.ParseInt32(&r.ReturnCode)
	buf.loadError("failed to parse return code")
}

func (buf *RespBuffer) parseOAUTH2ErrString(r *ResponseOAUTH2) {
	if buf.err != nil {
		return
	}
	buf.buffer.ParseString(&r.ErrorString)
	buf.loadError("failed to parse error string")
}

func (buf *RespBuffer) parseOAUTH2ClientID(r *ResponseOAUTH2) {
	if buf.err != nil {
		return
	}
	buf.buffer.ParseString(&r.CliendID)
	buf.loadError("failed to parse client id")
}

func (buf *RespBuffer) parseOAUTH2ClientType(r *ResponseOAUTH2) {
	if buf.err != nil {
		return
	}
	buf.buffer.ParseInt32(&r.ClientType)
	buf.loadError("failed to parse client type")
}

func (buf *RespBuffer) parseOAUTH2Username(r *ResponseOAUTH2) {
	if buf.err != nil {
		return
	}
	buf.buffer.ParseString(&r.Username)
	buf.loadError("failed to parse username")
}

func (buf *RespBuffer) parseOAUTH2ExpiresInInfo(r *ResponseOAUTH2) {
	if buf.err != nil {
		return
	}
	buf.buffer.ParseInt32(&r.ExpiresIn)
	buf.loadError("failed to parse expires_in data")
}

func (buf *RespBuffer) parseOAUTH2UserID(r *ResponseOAUTH2) {
	if buf.err != nil {
		return
	}
	buf.buffer.ParseInt64(&r.UserID)
	buf.loadError("failed to parse user id")
}
