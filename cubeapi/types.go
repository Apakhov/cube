package cubeapi

// Header represents header of response and request
type Header struct {
	SvcID      int32
	BodyLength int32
	RequestID  int32
}

const headerLen = 12

const int8Len = 1
const int32Len = 4
const int64Len = 8

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
		msg: "Not enough data",
	}
	// ErrIncorrectData can't parse due to incorrect data
	ErrIncorrectData = &Error{
		msg: "Incorrect data",
	}
	// ErrStringTooLong string is too long to write
	ErrStringTooLong = &Error{
		msg: "String is too long",
	}
	// ErrBadWritingPos can't write on this position
	ErrBadWritingPos = &Error{
		msg: "Can't write to this position",
	}
	// ErrIncorrectBodyLen incorrect body length
	ErrIncorrectBodyLen = &Error{
		msg: "Incorrect body length",
	}
	// ErrIncorrectSVCID incorrect svc id
	ErrIncorrectSVCID = &Error{
		msg: "Incorrect svc id",
	}
)
