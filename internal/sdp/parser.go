package sdp

import (
	"bytes"

	"github.com/indigo-web/utils/uf"
)

type Parser struct{}

func NewParser() Parser {
	return Parser{}
}

func (Parser) Parse(data []byte) (desc Description, err error) {
	session := Session{}
	var value string

	for len(data) > 0 {
		if len(data) < 2 {
			return desc, ErrIncompleteData
		}

		if data[1] != '=' {
			return desc, ErrBadSyntax
		}

		if data[0] == 'm' {
			// media starts here
			break
		}

		key := data[0]
		value, data = parseValue(data[2:])

		switch key {
		case 'v':
			session.Protocol = value
		case 'o':
			session.Originator, err = session.Originator.Parse(value)
			if err != nil {
				return desc, err
			}
		case 's':
			session.Name = value
		case 'i':
			session.Info = value
		case 'u':
			session.URI = value
		case 'e':
			session.Email = value
		case 'p':
			session.Phone = value
		case 'c':
			connInfo, err := ConnectionInfo{}.Parse(value)
			if err != nil {
				return desc, err
			}

			session.ConnectionInfo = append(session.ConnectionInfo, connInfo)
		case 'b':
			session.BandwidthInfo = append(session.BandwidthInfo, value)
		case 'z':
			session.TimeZoneAdjustments = append(session.TimeZoneAdjustments, value)
		case 'k':
			session.EncryptionKey = value
		case 'a':
			session.Attributes = append(session.Attributes, value)
		default:
			return desc, ErrUnrecognizedKey
		}
	}

	desc.Session = session

	if len(data) == 0 {
		return desc, nil
	}

	var media Media

	for len(data) > 0 {
		if len(data) < 2 {
			return desc, ErrIncompleteData
		}

		if data[1] != '=' {
			return desc, ErrBadSyntax
		}

		if data[0] == 'm' && media.Name != "" {
			// the next media block description has begun
			desc.Media = append(desc.Media, media)
			media = Media{}
			continue
		}

		key := data[0]
		value, data = parseValue(data[2:])

		switch key {
		case 'm':
			media.Name = value
		case 'i':
			media.Title = value
		case 'c':
			connInfo, err := ConnectionInfo{}.Parse(value)
			if err != nil {
				return desc, err
			}

			media.ConnectionInfo = append(media.ConnectionInfo, connInfo)
		case 'b':
			media.BandwidthInfo = append(media.BandwidthInfo, value)
		case 'k':
			media.EncryptionKey = value
		case 'a':
			media.Attributes = append(media.Attributes, value)
		default:
			return desc, ErrUnrecognizedKey
		}
	}

	desc.Media = append(desc.Media, media)

	return desc, nil
}

func parseValue(data []byte) (value string, rest []byte) {
	rest = data
	lf := bytes.IndexByte(data, '\n')
	if lf >= 0 {
		rest, data = data[lf+1:], data[:lf]

		if data[len(data)-1] == '\r' {
			data = data[:len(data)-1]
		}
	}

	return uf.B2S(data), rest
}
