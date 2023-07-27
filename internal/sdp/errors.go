package sdp

import "errors"

var (
	ErrUnrecognizedKey         = errors.New("met an unrecognized field identifier")
	ErrIncompleteData          = errors.New("provided data is incomplete")
	ErrBadSyntax               = errors.New("description is malformed")
	ErrUnknownNetType          = errors.New("received unsupported o=<nettype> value")
	ErrUnknownAddrType         = errors.New("received unsupported o=<addrtype> value")
	ErrUnknownEncryptionMethod = errors.New("received unknown encryption method")
)
