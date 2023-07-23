package sip

type parserState int

const (
	eMethod parserState = iota + 1
	eUriScheme
	eUriUser
	eUriUserD1
	eUriUserD2
	eUriPassword
	eUriPasswordD1
	eUriPasswordD2
	eUriHost
	eUriPort
	eParamsKey
	eParamsKeyD1
	eParamsKeyD2
	eParamsValue
	eParamsValueD1
	eParamsValueD2
	eProto
	eS
	eSI
	eSIP
	eProtoVersion
	eProtoCR
	eProtoCRLF
	eProtoCRLFCR
	eHeaderKey
	eHeaderColon
	eContentLength
	eContentLengthCR
	eContentLengthCRLF
	eContentLengthCRLFCR
	eHeaderValue
	eHeaderValueCR
	eHeaderValueCRLF
	eHeaderValueCRLFCR
	eBody
)
