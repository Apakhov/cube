package oauth2

import (
	"github.com/Apakhov/cube/cubeapi"
	"github.com/pkg/errors"
)

type SendBuffer struct {
	buffer *cubeapi.SendBuffer
}

func (buf *SendBuffer) Bytes() []byte {
	return buf.buffer.Bytes()
}

func CreateOAUTH2Request(token, scope string) (*SendBuffer, error) {
	buf := &SendBuffer{cubeapi.CreateSendBuffer()}
	bodyLen, err := buf.writeOAUTH2Body(token, scope)
	if err != nil {
		err = errors.Wrap(switchError(err), "failed to write request body")
		return nil, err
	}

	buf.buffer.WriteHeader(cubeOAUTH2SvcID, bodyLen)
	return buf, nil
}

func (buf *SendBuffer) writeOAUTH2Body(token, scope string) (bodyLen int32, err error) {
	headerLen := buf.buffer.Len()
	buf.buffer.WriteInt32(0)

	err = buf.buffer.WriteString(token)
	if err != nil {
		err = errors.Wrap(switchError(err), "can't write to buffer")
		return
	}
	err = buf.buffer.WriteString(scope)
	if err != nil {
		err = errors.Wrap(switchError(err), "can't write to buffer")
		return
	}
	bodyLen = int32(buf.buffer.Len() - headerLen)
	buf.buffer.WriteInt32OnPos(cubeOAUTH2SvcMSG, headerLen) //ignoring error since we know position exists
	return
}
