package parser

import (
	"fmt"
	"strings"

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
	tempParamKey      string
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
	case eUriScheme:
		goto uriScheme
	case eUriUser:
		goto uriUser
	case eUriUserD1:
		goto uriUserD1
	case eUriUserD2:
		goto uriUserD2
	case eUriPassword:
		goto uriPassword
	case eUriPasswordD1:
		goto uriPasswordD1
	case eUriPasswordD2:
		goto uriPasswordD2
	case eUriHost:
		goto uriHost
	case eUriPort:
		goto uriPort
	case eParamsKey:
		goto paramsKey
	case eParamsKeyD1:
		goto paramsKeyD1
	case eParamsKeyD2:
		goto paramsKeyD2
	case eParamsValue:
		goto paramsValue
	case eParamsValueD1:
		goto paramsValueD1
	case eParamsValueD2:
		goto paramsValueD2
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
			p.request.Method = uf.B2S(p.requestLineArena.Finish())
			data = data[i+1:]
			p.counter = 0
			p.state = eUriScheme
			goto uriScheme
		default:
			if !p.requestLineArena.Append(data[i]) {
				return true, status.ErrURITooLong
			}
		}
	}

	return false, nil

uriScheme:
	for i := range data {
		if data[i] == ':' {
			p.request.URI.Scheme = uf.B2S(data[:i])
			data = data[i+1:]
			p.state = eUriUser
			goto uriUser
		}
	}

	return false, nil

uriUser:
	for i := range data {
		switch data[i] {
		case '@':
			p.request.URI.User = uf.B2S(p.requestLineArena.Finish())
			data = data[i+1:]
			p.state = eUriHost
			goto uriHost
		case ':':
			p.request.URI.User = uf.B2S(p.requestLineArena.Finish())
			data = data[i+1:]
			p.state = eUriPassword
			goto uriPassword
		case '%':
			data = data[i+1:]
			p.state = eUriUserD1
			goto uriUserD1
		default:
			if !p.requestLineArena.Append(data[i]) {
				return true, status.ErrURITooLong
			}
		}
	}

	return false, nil

uriUserD1:
	if len(data) == 0 {
		return false, nil
	}

	if !isHex(data[0]) {
		return true, status.ErrURIDecoding
	}

	p.urlEncodedChar = unHex(data[0]) << 4
	data = data[1:]
	p.state = eUriUserD2
	goto uriUserD2

uriUserD2:
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
	p.state = eUriUser
	goto uriUser

uriPassword:
	for i := range data {
		switch data[i] {
		case '@':
			p.request.URI.Password = uf.B2S(p.requestLineArena.Finish())
			data = data[i+1:]
			p.state = eUriHost
			goto uriHost
		case '%':
			data = data[i+1:]
			p.state = eUriPasswordD1
			goto uriPasswordD1
		default:
			if !p.requestLineArena.Append(data[i]) {
				return true, status.ErrURITooLong
			}
		}
	}

	return false, nil

uriPasswordD1:
	if len(data) == 0 {
		return false, nil
	}

	if !isHex(data[0]) {
		return true, status.ErrURIDecoding
	}

	p.urlEncodedChar = unHex(data[0]) << 4
	data = data[1:]
	p.state = eUriPasswordD2
	goto uriPasswordD2

uriPasswordD2:
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
	p.state = eUriPassword
	goto uriPassword

uriHost:
	for i := range data {
		switch data[i] {
		case ' ':
			p.request.URI.Host = uf.B2S(data[:i])
			data = data[i+1:]
			p.state = eProto
			goto proto
		case ':':
			p.request.URI.Host = uf.B2S(data[:i])
			data = data[i+1:]
			p.state = eUriPort
			goto uriPort
		case ';':
			p.request.URI.Host = uf.B2S(data[:i])
			data = data[i+1:]
			p.state = eParamsKey
			goto paramsKey
		}
	}

	return false, nil

uriPort:
	for i := range data {
		switch data[i] {
		case ' ':
			p.request.URI.Port = p.counter
			data = data[i+1:]
			p.state = eProto
			goto proto
		case ';':
			p.request.URI.Port = p.counter
			data = data[i+1:]
			p.state = eParamsKey
			goto paramsKey
		default:
			if data[i] < '0' || data[i] > '9' {
				return true, status.ErrBadRequest
			}

			p.counter = p.counter*10 + int(data[i]-'0')
		}
	}

	return false, nil

paramsKey:
	for i := range data {
		switch data[i] {
		case '=':
			p.tempParamKey = uf.B2S(p.requestLineArena.Finish())
			data = data[i+1:]
			p.state = eParamsValue
			goto paramsValue
		case '%':
			data = data[i+1:]
			p.state = eParamsKeyD1
			goto paramsKeyD1
		case '+':
			if !p.requestLineArena.Append(' ') {
				return true, status.ErrURITooLong
			}
		default:
			if !p.requestLineArena.Append(data[i]) {
				return true, status.ErrURITooLong
			}
		}
	}

	return false, nil

paramsKeyD1:
	if len(data) == 0 {
		return false, nil
	}

	if !isHex(data[0]) {
		return true, status.ErrURIDecoding
	}

	p.urlEncodedChar = unHex(data[0]) << 4
	data = data[1:]
	p.state = eParamsKeyD2
	goto paramsKeyD2

paramsKeyD2:
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
	p.state = eParamsKey
	goto paramsKey

paramsValue:
	for i := range data {
		switch data[i] {
		case ';':
			p.request.URI.Params.Add(p.tempParamKey, uf.B2S(p.requestLineArena.Finish()))
			data = data[i+1:]
			p.state = eParamsKey
			goto paramsKey
		case ' ':
			p.request.URI.Params.Add(p.tempParamKey, uf.B2S(p.requestLineArena.Finish()))
			data = data[i+1:]
			p.state = eProto
			goto proto
		case '%':
			data = data[i+1:]
			p.state = eParamsValueD1
			goto paramsValueD1
		default:
			if !p.requestLineArena.Append(data[i]) {
				return true, status.ErrURITooLong
			}
		}
	}

	return false, nil

paramsValueD1:
	if len(data) == 0 {
		return false, nil
	}

	if !isHex(data[0]) {
		return true, status.ErrURIDecoding
	}

	p.urlEncodedChar = unHex(data[0]) << 4
	data = data[1:]
	p.state = eParamsValueD2
	goto paramsValueD2

paramsValueD2:
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
	p.state = eParamsValue
	goto paramsValue

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

	p.request.Proto = sip.Protocol(uf.B2S(p.requestLineArena.Finish()))

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

			if strings.EqualFold(p.headerKey, "content-length") {
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
