package oauth2

import (
	"github.com/Apakhov/cube/cubeapi"
	"github.com/pkg/errors"
)

type ParseError struct {
	msg string
}

func (e *ParseError) Error() string {
	return e.msg
}

var (
	// ErrNotEnoughData not enough data to parse
	ErrNotEnoughData = &ParseError{
		msg: "oauth2: Not enough data",
	}
	// ErrIncorrectData can't parse due to incorrect data
	ErrIncorrectData = &ParseError{
		msg: "oauth2: Incorrect data",
	}
	// ErrStringTooLong string is too long to write
	ErrStringTooLong = &ParseError{
		msg: "oauth2: String is too long",
	}
	// ErrBadWritingPos can't write on this position
	ErrBadWritingPos = &ParseError{
		msg: "oauth2: Can't write to this position",
	}
	// ErrIncorrectBodyLen incorrect body length
	ErrIncorrectBodyLen = &ParseError{
		msg: "oauth2: Incorrect body length",
	}
	// ErrIncorrectSVCID incorrect svc id
	ErrIncorrectSVCID = &ParseError{
		msg: "oauth2: Incorrect svc id",
	}
	// ErrUndefined error is not supported
	ErrUndefined = &ParseError{
		msg: "oauth2: error is not supported",
	}
)

func switchError(e error) *ParseError {
	switch c := errors.Cause(e).(type) {
	case *cubeapi.ParseError:
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
		case cubeapi.ErrIncorrectSVCID:
			return ErrIncorrectSVCID
		default:
			return ErrUndefined
		}
	case *ParseError:
		return c
	default:
		return &ParseError{
			msg: c.Error(),
		}
	}
}
