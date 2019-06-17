package cubeapi

type Header struct {
	SvcID      int32
	BodyLength int32
	RequestID  int32
}

const headerLen = 12

const int32Len = 4
const int64Len = 8

type ParseError struct {
	msg string
}

func (e *ParseError) Error() string {
	return e.msg
}

var (
	// ErrNotEnoughData not enough data to parse
	ErrNotEnoughData = &ParseError{
		msg: "Not enough data",
	}
	// ErrIncorrectData can't parse due to incorrect data
	ErrIncorrectData = &ParseError{
		msg: "Incorrect data",
	}
	// ErrStringTooLong string is too long to write
	ErrStringTooLong = &ParseError{
		msg: "String is too long",
	}
	// ErrBadWritingPos can't write on this position
	ErrBadWritingPos = &ParseError{
		msg: "Can't write to this position",
	}
	// ErrIncorrectBodyLen incorrect body length
	ErrIncorrectBodyLen = &ParseError{
		msg: "Incorrect body length",
	}
	// ErrIncorrectSVCID incorrect svc id
	ErrIncorrectSVCID = &ParseError{
		msg: "Incorrect svc id",
	}
)
