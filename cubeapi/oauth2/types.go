package oauth2

import (
	"fmt"
)

// ResponseOAUTH2 represents oauth2 response
type ResponseOAUTH2 struct {
	ReturnCode  int32
	CliendID    string
	ClientType  int32
	Username    string
	ExpiresIn   int32
	UserID      int64
	ErrorString string
}

const cubeOAUTH2SvcID = int32(0x00000002)
const cubeOAUTH2SvcMSG = int32(0x00000001)

// codes of errors
const (
	CubeOAUTH2ErrCodeOK = int32(iota)
	CubeOAUTH2ErrCodeTokenNotFound
	CubeOAUTH2ErrCodeDBError
	CubeOAUTH2ErrCodeUnknownMSG
	CubeOAUTH2ErrCodeBadPacket
	CubeOAUTH2ErrCodeBadClient
	CubeOAUTH2ErrCodeBadScope
)

// error strings
const (
	CubeOAUTH2ErrStringOK            = "CUBE_OAUTH2_ERR_OK"
	CubeOAUTH2ErrStringTokenNotFound = "CUBE_OAUTH2_ERR_TOKEN_NOT_FOUND"
	CubeOAUTH2ErrStringDBError       = "CUBE_OAUTH2_ERR_DB_ERROR"
	CubeOAUTH2ErrStringUnknownMSG    = "CUBE_OAUTH2_ERR_UNKNOWN_MSG"
	CubeOAUTH2ErrStringBadPacket     = "CUBE_OAUTH2_ERR_BAD_PACKET"
	CubeOAUTH2ErrStringBadClient     = "CUBE_OAUTH2_ERR_BAD_CLIENT"
	CubeOAUTH2ErrStringBadScope      = "CUBE_OAUTH2_ERR_BAD_SCOPE"
)

// errors description
const (
	CubeOAUTH2ErrDescrOK            = "---------------"
	CubeOAUTH2ErrDescrTokenNotFound = "token not found"
	CubeOAUTH2ErrDescrDBError       = "db error"
	CubeOAUTH2ErrDescrUnknownMSG    = "unknown svc message type"
	CubeOAUTH2ErrDescrBadPacket     = "bad packet"
	CubeOAUTH2ErrDescrBadClient     = "bad client"
	CubeOAUTH2ErrDescrBadScope      = "bad scope"
)

func errInfoByCode(c int32) (string, string) {
	switch c {
	case CubeOAUTH2ErrCodeOK:
		return CubeOAUTH2ErrDescrOK, CubeOAUTH2ErrStringOK
	case CubeOAUTH2ErrCodeTokenNotFound:
		return CubeOAUTH2ErrDescrTokenNotFound, CubeOAUTH2ErrStringTokenNotFound
	case CubeOAUTH2ErrCodeDBError:
		return CubeOAUTH2ErrDescrDBError, CubeOAUTH2ErrStringDBError
	case CubeOAUTH2ErrCodeUnknownMSG:
		return CubeOAUTH2ErrDescrUnknownMSG, CubeOAUTH2ErrStringUnknownMSG
	case CubeOAUTH2ErrCodeBadPacket:
		return CubeOAUTH2ErrDescrBadPacket, CubeOAUTH2ErrStringBadPacket
	case CubeOAUTH2ErrCodeBadClient:
		return CubeOAUTH2ErrDescrBadClient, CubeOAUTH2ErrStringBadClient
	case CubeOAUTH2ErrCodeBadScope:
		return CubeOAUTH2ErrDescrBadScope, CubeOAUTH2ErrStringBadScope

	default:
		return "unknown error code", "unknown error code"
	}
}

func (r *ResponseOAUTH2) String() string {
	if r.ReturnCode == CubeOAUTH2ErrCodeOK {
		errStr, errDescr := errInfoByCode(r.ReturnCode)
		return fmt.Sprintf(`error: %s
message: %s`, errStr, errDescr)
	}

	return fmt.Sprintf(`client_id: %s
client_type: %d
expires_in: %d
user_id: %d
username: %s`, r.CliendID, r.ClientType, r.ExpiresIn, r.UserID, r.Username)
}
