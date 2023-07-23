package sdp

import "errors"

var (
	ErrUnrecognizedKey = errors.New("met an unrecognized field identifier")
	ErrInvalidCRLF     = errors.New("wanted CRLF, got CR with unknown following char")
	ErrIncompleteData  = errors.New("provided data is incomplete")
	ErrBadSyntax       = errors.New("description is malformed")
)
