package sdp

import (
	"strings"
)

type NetType string

const (
	IN NetType = "IN"
)

type AddrType string

const (
	IP4 AddrType = "IP4"
	IP6 AddrType = "IP6"
)

type Origin struct {
	Username       string
	SessId         string // TODO: these may be ordinary integers
	SessVersion    string // TODO: these may be ordinary integers
	NetType        NetType
	AddrType       AddrType
	UnicastAddress string
}

func (o Origin) Parse(value string) (Origin, error) {
	sp := strings.IndexByte(value, ' ')
	if sp == -1 {
		return o, ErrBadSyntax
	}

	o.Username, value = value[:sp], value[sp+1:]

	sp = strings.IndexByte(value, ' ')
	if sp == -1 {
		return o, ErrBadSyntax
	}

	o.SessId, value = value[:sp], value[sp+1:]

	sp = strings.IndexByte(value, ' ')
	if sp == -1 {
		return o, ErrBadSyntax
	}

	o.SessVersion, value = value[:sp], value[sp+1:]

	sp = strings.IndexByte(value, ' ')
	if sp == -1 {
		return o, ErrBadSyntax
	}

	switch nettype := value[:sp]; nettype {
	case "IN":
		o.NetType = IN
	default:
		return o, ErrUnknownNetType
	}

	value = value[sp+1:]

	sp = strings.IndexByte(value, ' ')
	if sp == -1 {
		return o, ErrBadSyntax
	}

	switch addrtype := value[:sp]; addrtype {
	case "IP4":
		o.AddrType = IP4
	case "IP6":
		o.AddrType = IP6
	default:
		return o, ErrUnknownAddrType
	}

	value = value[sp+1:]

	if len(value) == 0 {
		return o, ErrBadSyntax
	}

	o.UnicastAddress = value

	return o, nil
}
