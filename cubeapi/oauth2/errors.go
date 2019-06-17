package oauth2

import (
	"github.com/Apakhov/cube/cubeapi"
	"github.com/pkg/errors"
)

// Error errors
type Error struct {
	msg string
}

// Error imlements error interface
func (e *Error) Error() string {
	return e.msg
}

var (
	// ErrNotEnoughData not enough data to parse
	ErrNotEnoughData = &Error{
		msg: "oauth2: Not enough data",
	}
	// ErrIncorrectData can't parse due to incorrect data
	ErrIncorrectData = &Error{
		msg: "oauth2: Incorrect data",
	}
	// ErrStringTooLong string is too long to write
	ErrStringTooLong = &Error{
		msg: "oauth2: String is too long",
	}
	// ErrBadWritingPos can't write on this position
	ErrBadWritingPos = &Error{
		msg: "oauth2: Can't write to this position",
	}
	// ErrIncorrectBodyLen incorrect body length
	ErrIncorrectBodyLen = &Error{
		msg: "oauth2: Incorrect body length",
	}
	// ErrIncorrectLen incorrect length of element
	ErrIncorrectLen = &Error{
		msg: "oauth2: Incorrect length of element",
	}
	// ErrIncorrectSVCID incorrect svc id
	ErrIncorrectSVCID = &Error{
		msg: "oauth2: Incorrect svc id",
	}
	// ErrUndefined error is not supported
	ErrUndefined = &Error{
		msg: "oauth2: error is not supported",
	}
)

func switchError(e error) *Error {
	switch c := errors.Cause(e).(type) {
	case *cubeapi.Error:
		switch c {
		case cubeapi.ErrNotEnoughData:
			return ErrNotEnoughData
		case cubeapi.ErrIncorrectData:
			return ErrIncorrectData
		case cubeapi.ErrStringTooLong:
			return ErrStringTooLong
		case cubeapi.ErrBadWritingPos:
			return ErrBadWritingPos
		case cubeapi.ErrIncorrectBodyLen:
			return ErrIncorrectBodyLen
		case cubeapi.ErrIncorrectLen:
			return ErrIncorrectLen
		case cubeapi.ErrIncorrectSVCID:
			return ErrIncorrectSVCID
		default:
			return ErrUndefined
		}
	case *Error:
		return c
	default:
		return &Error{
			msg: c.Error(),
		}
	}
}
