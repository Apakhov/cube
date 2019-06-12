package oauth2

import (
	"bytes"
	"encoding/binary"

	"github.com/pkg/errors"
)

const headerLen = 12

func createOAUTH2Request(token, scope string) ([]byte, error) {
	buf := createResponse()
	buf, bodyLen, err := writeOAUTH2Body(token, scope, buf)
	if err != nil {
		err = errors.Wrap(err, "failed to write request body")
		return nil, err
	}

	buf = writeHeader(cubeOAUTH2SvcID, bodyLen, buf)
	return buf, nil
}

const cubeOAUTH2SvcID = int32(0x00000002)

func createResponse() []byte {
	return make([]byte, headerLen, headerLen)
}

func writeHeader(svcID int32, bodyLen int32, buf []byte) []byte {
	binary.LittleEndian.PutUint32(buf[0:4], uint32(svcID))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(bodyLen))
	binary.LittleEndian.PutUint32(buf[8:12], 0x00000000)
	return buf
}

const cubeOAUTH2SvcMSG = int32(0x00000001)

func writeOAUTH2Body(token, scope string, buf []byte) (res []byte, bodyLen int32, err error) {
	buffer := bytes.NewBuffer(buf)

	var cur int
	var bLen int
	cur, err = buffer.Write([]byte{0, 0, 0, 0}) //reserve space for svc_msg
	if err != nil {
		err = errors.Wrap(err, "can't write to buffer")
		return
	}
	bLen += cur
	cur, err = buffer.WriteString(token)
	if err != nil {
		err = errors.Wrap(err, "can't write to buffer")
		return
	}
	bLen += cur
	cur, err = buffer.WriteString(scope)
	if err != nil {
		err = errors.Wrap(err, "can't write to buffer")
		return
	}
	bLen += cur
	res = buffer.Bytes()

	binary.LittleEndian.PutUint32(res[headerLen:headerLen+4], uint32(cubeOAUTH2SvcMSG))
	bodyLen = int32(bLen)
	return
}
