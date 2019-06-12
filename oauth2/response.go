package oauth2

import (
	"bytes"

	"github.com/pkg/errors"
)

type responseOAUTH2 struct {
	returnCode  int32
	cliendID    string
	clientType  int32
	username    string
	expiresIn   int32
	userID      int64
	errorString string
}

func parseOAUTH2Resp(buffer *bytes.Buffer) (r responseOAUTH2, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "failed to parse OAUTH2 response")
		}
	}()
	h, err := cubeapi.parseHeader(buffer)
	if err != nil {
		return
	}
	if int(h.bodyLength) != buffer.Len() {
		err = errors.New("incorrect body length")
		return
	}
	r, err = parseOAUTH2RespBody(buffer)
	return
}

func parseOAUTH2RespBody(buffer *bytes.Buffer) (r responseOAUTH2, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "failed to parse OAUTH2 response body")
		}
	}()
	err = parseOAUTH2RespReturnCode(&r, buffer)
	if err != nil {
		return
	}
	if r.returnCode != 0x00000000 {
		err = parseOAUTH2ErrString(&r, buffer)
		return
	}
	err = parseOAUTH2ClientID(&r, buffer)
	if err != nil {
		return
	}
	err = parseOAUTH2ClientType(&r, buffer)
	if err != nil {
		return
	}
	err = parseOAUTH2Username(&r, buffer)
	if err != nil {
		return
	}
	err = parseOAUTH2ExpiresInInfo(&r, buffer)
	if err != nil {
		return
	}
	err = parseOAUTH2UserID(&r, buffer)
	return
}

func parseOAUTH2RespReturnCode(r *responseOAUTH2, buffer *bytes.Buffer) error {
	err := cubeapi.parseInt32(&r.returnCode, buffer)
	if err != nil {
		return errors.Wrap(err, "failed to parse return code")
	}
	return nil
}

func parseOAUTH2ErrString(r *responseOAUTH2, buffer *bytes.Buffer) error {
	err := cubeapi.parseString(&r.errorString, buffer)
	if err != nil {
		return errors.Wrap(err, "failed to parse error string")
	}
	return nil
}

func parseOAUTH2ClientID(r *responseOAUTH2, buffer *bytes.Buffer) error {
	err := cubeapi.parseString(&r.cliendID, buffer)
	if err != nil {
		return errors.Wrap(err, "failed to parse client id")
	}
	return nil
}

func parseOAUTH2ClientType(r *responseOAUTH2, buffer *bytes.Buffer) error {
	err := cubeapi.parseInt32(&r.clientType, buffer)
	if err != nil {
		return errors.Wrap(err, "failed to parse client type")
	}
	return nil
}

func parseOAUTH2Username(r *responseOAUTH2, buffer *bytes.Buffer) error {
	err := cubeapi.parseString(&r.username, buffer)
	if err != nil {
		return errors.Wrap(err, "failed to parse username")
	}
	return nil
}

func parseOAUTH2ExpiresInInfo(r *responseOAUTH2, buffer *bytes.Buffer) error {
	err := cubeapi.parseInt32(&r.expiresIn, buffer)
	if err != nil {
		return errors.Wrap(err, "failed to parse expires_in data")
	}
	return nil
}

func parseOAUTH2UserID(r *responseOAUTH2, buffer *bytes.Buffer) error {
	err := cubeapi.parseInt64(&r.userID, buffer)
	if err != nil {
		return errors.Wrap(err, "failed to parse user id")
	}
	return nil
}
