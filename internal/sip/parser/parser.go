package parser

import (
	"fmt"

	"github.com/gokiki/sip-server/internal/sip"
	"github.com/gokiki/sip-server/internal/sip/status"
	"github.com/gokiki/sip-server/settings"
	"github.com/indigo-web/utils/arena"
	"github.com/indigo-web/utils/pool"
	"github.com/indigo-web/utils/uf"
)

var _ sip.Parser = &Parser{}

type Parser struct {
	request           *sip.Request
	headerKey         string
	requestLineArena  arena.Arena[byte]
	headersValuesPool pool.ObjectPool[[]string]
	headerKeyArena    arena.Arena[byte]
	headerValueArena  arena.Arena[byte]
	settings          settings.Settings
	// generic multi-purpose counter
	counter        int
	contentLength  int
	headerSize     int
	urlEncodedChar uint8
	bodyBuff       []byte
	state          parserState
}

func NewParser(
	request *sip.Request, keyArena, valArena, requestLineArena arena.Arena[byte],
	valuesPool pool.ObjectPool[[]string], s settings.Settings,
) *Parser {
	return &Parser{
		state:             eMethod,
		request:           request,
		settings:          s,
		requestLineArena:  requestLineArena,
		headerKeyArena:    keyArena,
		headerValueArena:  valArena,
		headersValuesPool: valuesPool,
	}
}

func (p *Parser) Parse(data []byte) (done bool, err error) {
	var value string
	_ = *p.request
	requestHeaders := p.request.Headers

	switch p.state {
	case eMethod:
		goto method
	case eUri:
		goto uri
	case eUriDecode1Char:
		goto uriDecode1Char
	case eUriDecode2Char:
		goto uriDecode2Char
	case eParams:
		goto params
	case eParamsDecode1Char:
		goto queryDecode1Char
	case eParamsDecode2Char:
		goto queryDecode2Char
	case eProto:
		goto proto
	case eS:
		goto protoS
	case eSI:
		goto protoSI
	case eSIP:
		goto protoSIP
	case eProtoVersion:
		goto protoVersion
	case eProtoCR:
		goto protoCR
	case eProtoCRLF:
		goto protoCRLF
	case eProtoCRLFCR:
		goto protoCRLFCR
	case eHeaderKey:
		goto headerKey
	case eHeaderColon:
		goto headerColon
	case eContentLength:
		goto contentLength
	case eContentLengthCR:
		goto contentLengthCR
	case eContentLengthCRLF:
		goto contentLengthCRLF
	case eContentLengthCRLFCR:
		goto contentLengthCRLFCR
	case eHeaderValue:
		goto headerValue
	case eHeaderValueCR:
		goto headerValueCR
	case eHeaderValueCRLF:
		goto headerValueCRLF
	case eHeaderValueCRLFCR:
		goto headerValueCRLFCR
	default:
		panic(fmt.Sprintf("BUG: unexpected state: %v", p.state))
	}

method:
	for i := range data {
		switch data[i] {
		case '\r', '\n':
			return true, status.ErrBadRequest
		case ' ':
			if p.counter == 0 {
				return true, status.ErrBadRequest
			}

			if !p.requestLineArena.Append(data[:i]...) {
				// FIXME: sometimes this may be caused by misconfiguration. Error message
				//  must be a bit more informative, but how?
				return true, status.ErrMethodNotImplemented
			}
			p.request.Method = uf.B2S(p.requestLineArena.Finish())
			data = data[i+1:]
			p.counter = 0
			p.state = eUri
			goto uri
		default:
			p.counter++

			if p.counter > p.settings.RequestLine.MaxMethodLength {
				return true, status.ErrBadRequest
			}
		}
	}

	if !p.requestLineArena.Append(data...) {
		return true, status.ErrMethodNotImplemented
	}

	return false, nil

uri:
	// FIXME: parse here not like http path, but sip address

	for i := range data {
		switch data[i] {
		case ' ':
			if p.begin == p.pointer {
				return true, status.ErrBadRequest
			}

			p.request.URI = uf.B2S(p.requestLineBuff[p.begin:p.pointer])
			data = data[i+1:]
			p.state = eProto
			p.begin = p.pointer
			goto proto
		case '%':
			data = data[i+1:]
			p.state = eUriDecode1Char
			goto uriDecode1Char
		case '?':
			p.request.Path.String = uf.B2S(p.requestLineBuff[p.begin:p.pointer])
			if len(p.request.Path.String) == 0 {
				p.request.Path.String = "/"
			}

			p.begin = p.pointer
			data = data[i+1:]
			p.state = eQuery
			goto query
		case '\x00', '\n', '\r', '\t', '\b', '\a', '\v', '\f':
			// request path MUST NOT include any non-printable characters
			return true, status.ErrBadRequest
		default:
			if p.pointer >= len(p.requestLineBuff) {
				return true, status.ErrURITooLong
			}

			p.requestLineBuff[p.pointer] = data[i]
			p.pointer++
		}
	}

	return false, nil

uriDecode1Char:
	if len(data) == 0 {
		return false, nil
	}

	if !isHex(data[0]) {
		return true, status.ErrURIDecoding
	}

	p.urlEncodedChar = unHex(data[0]) << 4
	data = data[1:]
	p.state = eUriDecode2Char
	goto uriDecode2Char

uriDecode2Char:
	if len(data) == 0 {
		return false, nil
	}

	if !isHex(data[0]) {
		return true, status.ErrURIDecoding
	}

	if !p.requestLineArena.Append(p.urlEncodedChar | unHex(data[0])) {
		return true, status.ErrURITooLong
	}

	data = data[1:]
	p.state = eUri
	goto uri

params:
	// FIXME: we need to parse all the params in-place. Replace this with a better
	//  state machine

	for i := range data {
		switch data[i] {
		case ' ':
			p.request.Path.Query.Set(p.requestLineBuff[p.begin:p.pointer])
			data = data[i+1:]
			p.state = eProto
			goto proto
		case '%':
			data = data[i+1:]
			p.state = eParamsDecode1Char
			goto paramsDecode1Char
		case '+':
			if p.pointer >= len(p.requestLineBuff) {
				return true, status.ErrURITooLong
			}

			if !p.requestLineArena.Append(' ') {
				return true, status.ErrURITooLong
			}
		case '\x00', '\n', '\r', '\t', '\b', '\a', '\v', '\f':
			return true, status.ErrBadRequest
		}
	}

	return false, nil

paramsDecode1Char:
	if len(data) == 0 {
		return false, nil
	}

	if !isHex(data[0]) {
		return true, status.ErrURIDecoding
	}

	p.urlEncodedChar = unHex(data[0]) << 4
	data = data[1:]
	p.state = eParamsDecode2Char
	goto paramsDecode2Char

paramsDecode2Char:
	if len(data) == 0 {
		return false, nil
	}

	if !isHex(data[0]) {
		return true, status.ErrURIDecoding
	}
	if !p.requestLineArena.Append(p.urlEncodedChar | unHex(data[0])) {
		return true, status.ErrURITooLong
	}

	data = data[1:]
	p.state = eParams
	goto params

proto:
	if len(data) == 0 {
		return false, nil
	}

	if data[0]|0x20 == 's' {
		if !p.requestLineArena.Append(data[0]) {
			return true, status.ErrURITooLong
		}

		data = data[1:]
		p.state = eS
		goto protoS
	}

	return false, status.ErrBadRequest

protoS:
	if len(data) == 0 {
		return false, nil
	}

	if data[0]|0x20 == 'i' {
		if !p.requestLineArena.Append(data[0]) {
			return true, status.ErrURITooLong
		}

		data = data[1:]
		p.state = eSI
		goto protoSI
	}

	return true, status.ErrUnsupportedProtocol

protoSI:
	if len(data) == 0 {
		return false, nil
	}

	if data[0]|0x20 == 'p' {
		if !p.requestLineArena.Append(data[0]) {
			return true, status.ErrURITooLong
		}

		data = data[1:]
		p.state = eSIP
		goto protoSIP
	}

	return true, status.ErrUnsupportedProtocol

protoSIP:
	if len(data) == 0 {
		return false, nil
	}

	if data[0] == '/' {
		if !p.requestLineArena.Append(data[0]) {
			return true, status.ErrURITooLong
		}

		data = data[1:]
		p.state = eProtoVersion
		goto protoVersion
	}

	return true, status.ErrUnsupportedProtocol

protoVersion:
	for i := range data {
		switch data[i] {
		case '\r':
			if !p.requestLineArena.Append(data[:i]...) {
				return true, status.ErrURITooLong
			}

			data = data[i+1:]
			p.state = eProtoCR
			goto protoCR
		case '\n':
			if !p.requestLineArena.Append(data[:i]...) {
				return true, status.ErrURITooLong
			}

			data = data[i+1:]
			p.state = eProtoCRLF
			goto protoCRLF
		}
	}

	if !p.requestLineArena.Append(data...) {
		return true, status.ErrURITooLong
	}

	return false, nil

protoCR:
	if len(data) == 0 {
		return false, nil
	}

	if data[0] != '\n' {
		return true, status.ErrBadRequest
	}

	data = data[1:]
	p.state = eProtoCRLF
	goto protoCRLF

protoCRLF:
	if len(data) == 0 {
		return false, nil
	}

	p.request.Proto = sip.Protocol(uf.B2S(p.requestLineBuff))

	switch data[0] {
	case '\r':
		data = data[1:]
		p.state = eProtoCRLFCR
		goto protoCRLFCR
	case '\n':
		if !p.request.HasBody() {
			return true, nil
		}

		data = data[1:]
		p.state = eBody
		goto body
	default:
		p.state = eHeaderKey
		goto headerKey
	}

protoCRLFCR:
	if len(data) == 0 {
		return false, nil
	}

	if data[0] == '\n' {
		return true, nil
	}

	return true, status.ErrBadRequest

headerKey:
	for i := range data {
		switch data[i] {
		case ':':
			p.counter++

			if p.counter > p.settings.Headers.MaxNumber {
				return true, status.ErrTooManyHeaders
			}

			if !p.headerKeyArena.Append(data[:i]...) {
				return true, status.ErrHeaderFieldsTooLarge
			}

			p.headerKey = uf.B2S(p.headerKeyArena.Finish())
			data = data[i+1:]

			if p.headerKey == "content-length" {
				p.state = eContentLength
				goto contentLength
			}

			p.state = eHeaderColon
			goto headerColon
		case '\r', '\n':
			return true, status.ErrBadRequest
		}
	}

	if !p.headerKeyArena.Append(data...) {
		return true, status.ErrHeaderFieldsTooLarge
	}

	return false, nil

headerColon:
	for i := range data {
		switch data[i] {
		case '\r', '\n':
			return true, status.ErrBadRequest
		case ' ':
		default:
			data = data[i:]
			p.state = eHeaderValue
			goto headerValue
		}
	}

	return false, nil

contentLength:
	for i := range data {
		switch char := data[i]; char {
		case ' ':
		case '\r':
			data = data[i+1:]
			p.state = eContentLengthCR
			goto contentLengthCR
		case '\n':
			data = data[i+1:]
			p.state = eContentLengthCRLF
			goto contentLengthCRLF
		default:
			if char < '0' || char > '9' {
				return true, status.ErrBadRequest
			}

			p.contentLength = p.contentLength*10 + int(char-'0')
		}
	}

	return false, nil

contentLengthCR:
	if len(data) == 0 {
		return false, nil
	}

	if data[0] == '\n' {
		data = data[1:]
		p.state = eContentLengthCRLF
		goto contentLengthCRLF
	}

	return true, status.ErrBadRequest

contentLengthCRLF:
	if len(data) == 0 {
		return false, nil
	}

	p.request.ContentLength = p.contentLength

	switch data[0] {
	case '\r':
		data = data[1:]
		p.state = eContentLengthCRLFCR
		goto contentLengthCRLFCR
	case '\n':
		if !p.request.HasBody() {
			return true, nil
		}

		data = data[1:]
		p.state = eBody
		goto body
	default:
		p.state = eHeaderKey
		goto headerKey
	}

contentLengthCRLFCR:
	if len(data) == 0 {
		return false, nil
	}

	if data[0] == '\n' {
		if !p.request.HasBody() {
			return true, nil
		}

		data = data[1:]
		p.state = eBody
		goto body
	}

	return true, status.ErrBadRequest

headerValue:
	for i := range data {
		switch data[i] {
		case '\r':
			if !p.headerValueArena.Append(data[:i]...) {
				return true, status.ErrHeaderFieldsTooLarge
			}

			data = data[i+1:]
			p.state = eHeaderValueCR
			goto headerValueCR
		case '\n':
			if !p.headerValueArena.Append(data[:i]...) {
				return true, status.ErrHeaderFieldsTooLarge
			}

			data = data[i+1:]
			p.state = eHeaderValueCRLF
			goto headerValueCRLF
		}
	}

	if !p.headerValueArena.Append(data...) {
		return true, status.ErrHeaderFieldsTooLarge
	}

	return false, nil

headerValueCR:
	if len(data) == 0 {
		return false, nil
	}

	if data[0] == '\n' {
		data = data[1:]
		p.state = eHeaderValueCRLF
		goto headerValueCRLF
	}

	return true, status.ErrBadRequest

headerValueCRLF:
	if len(data) == 0 {
		return false, nil
	}

	value = uf.B2S(p.headerValueArena.Finish())
	requestHeaders.Add(p.headerKey, value)

	switch data[0] {
	case '\n':
		if !p.request.HasBody() {
			return true, nil
		}

		data = data[1:]
		p.state = eBody
		goto body
	case '\r':
		data = data[1:]
		p.state = eHeaderValueCRLFCR
		goto headerValueCRLFCR
	default:
		p.state = eHeaderKey
		goto headerKey
	}

headerValueCRLFCR:
	if len(data) == 0 {
		return false, nil
	}

	if data[0] == '\n' {
		if !p.request.HasBody() {
			return true, nil
		}

		data = data[1:]
		p.state = eBody
		goto body
	}

	return true, status.ErrBadRequest

body:
	if len(data) < p.contentLength {
		p.bodyBuff = append(p.bodyBuff, data...)
		p.contentLength -= len(data)

		return false, nil
	}

	p.bodyBuff = append(p.bodyBuff, data[:p.contentLength]...)
	p.contentLength = 0
	p.request.Body = p.bodyBuff

	return true, nil
}

func (p *Parser) Release() {
	p.request.Headers.Clear()
	p.headerKeyArena.Clear()
	p.headerValueArena.Clear()
	p.requestLineArena.Clear()
	p.state = eMethod
}
