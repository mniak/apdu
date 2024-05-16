package apdu

import (
	"fmt"
)

const (
	ErrNoInformationGiven                       TrailerError = 0x6A00
	ErrIncorrectParametersInTheCommandDataField TrailerError = 0x6A80
	ErrFunctionNotSupported                     TrailerError = 0x6A81
	ErrFileOrApplicationNotFound                TrailerError = 0x6A82
	ErrRecordNotFound                           TrailerError = 0x6A83
	ErrNotEnoughMemorySpaceInTheFile            TrailerError = 0x6A84
	ErrNcInconsistentWithTLVStructure           TrailerError = 0x6A85
	ErrIncorrectParametersP1P2                  TrailerError = 0x6A86
	ErrNcInconsistentWithParametersP1P2         TrailerError = 0x6A87
	ErrReferencedDataOrReferenceDataNotFound    TrailerError = 0x6A88
	ErrFileAlreadyExists                        TrailerError = 0x6A89
	ErrDFNameAlreadyExists                      TrailerError = 0x6A8A
)

func (t Trailer) GetError() error {
	if t == 0x9000 {
		return nil
	}
	return TrailerError(t)
}

type TrailerError uint16

func (te TrailerError) Error() string {
	msg := Trailer(te).message()
	if msg != "" {
		return fmt.Sprintf("invalid trailer [%s]: %s", Trailer(te).code(), msg)
	} else {
		return fmt.Sprintf("invalid trailer [%s]", Trailer(te).code())
	}
}

func (te TrailerError) SW1() byte {
	return Trailer(te).SW1()
}

func (te TrailerError) SW2() byte {
	return Trailer(te).SW1()
}

func IsTrailerError(err error, val TrailerError) bool {
	te, is := AsTrailerError(err)
	return is && te == val
}

func AsTrailerError(err error) (TrailerError, bool) {
	if err == nil {
		return 0, false
	}

	te, is := err.(TrailerError)
	if is {
		return te, true
	}
	teptr, is := err.(*TrailerError)
	if is {
		return *teptr, true
	}

	unwrappable, can := err.(interface{ Unwrap(error) error })
	if can {
		internalErr := unwrappable.Unwrap(err)
		return AsTrailerError(internalErr)
	}

	return 0, false
}
