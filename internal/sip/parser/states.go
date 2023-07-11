package parser

type parserState int

const (
	eMethod parserState = iota + 1
	eUri
	eUriDecode1Char
	eUriDecode2Char
	eParams
	eParamsDecode1Char
	eParamsDecode2Char
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
