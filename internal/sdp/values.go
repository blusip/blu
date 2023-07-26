package sdp

import (
	"strconv"
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

type BandwidthType string

const (
	UnknownBWType BandwidthType = ""
	CT            BandwidthType = "CT"
	AS            BandwidthType = "AS"
)

type Bandwidth struct {
	Type  BandwidthType
	Value int
}

func (b Bandwidth) Parse(value string) (bandwidth Bandwidth, err error) {
	colon := strings.IndexByte(value, ':')
	if colon == -1 {
		return b, ErrBadSyntax
	}

	btype, bvalue := value[:colon], value[colon+1:]
	switch btype {
	case "CT":
		b.Type = CT
	case "AS":
		b.Type = AS
	default:
		b.Type = UnknownBWType
		return b, nil
	}

	b.Value, err = strconv.Atoi(bvalue)
	if err != nil {
		return b, err
	}

	return b, nil
}

type ConnectionInfo struct {
	NetType   NetType
	AddrType  AddrType
	Address   string
	TTL       int
	AddrRange int
}

func (c ConnectionInfo) Parse(value string) (info ConnectionInfo, err error) {
	sp := strings.IndexByte(value, ' ')
	if sp == -1 {
		return c, ErrBadSyntax
	}

	switch nettype := value[:sp]; nettype {
	case "IN":
		c.NetType = IN
	default:
		return c, ErrUnknownNetType
	}

	value = value[sp+1:]

	sp = strings.IndexByte(value, ' ')
	if sp == -1 {
		return c, ErrBadSyntax
	}

	switch addrtype := value[:sp]; addrtype {
	case "IP4":
		c.AddrType = IP4
	case "IP6":
		c.AddrType = IP6
	default:
		return c, ErrUnknownAddrType
	}

	value = value[sp+1:]

	if sp = strings.IndexByte(value, ' '); sp != -1 {
		value = value[:sp]
	}

	c.Address, c.TTL, c.AddrRange, err = parseAddress(value, c.AddrType)

	return c, err
}

// parseAddress parses the address provided in connection info field. It follows
// the following grammar:
//
//	(ip4addr | ip6addr) [ "/" ttl [ "/" addr-range ] ]
//
// If no TTL is provided, -1 will be returned.
// If no addr-range is provided, 0 will be returned.
//
// Actually this grammar doesn't follow the official one from the RFC. Those one is pretty confusing,
// as TTL there is optional, while address range isn't. But in case there is just one element
// after the address separated with a slash, then it isn't an addr-range, but TTL. In case you want
// to specify the addr-range, then you MUST also specify TTL.
func parseAddress(addr string, typ AddrType) (outAddr string, ttl, addrrange int, err error) {
	slash := strings.IndexByte(addr, '/')
	if slash == -1 {
		return addr, -1, 0, nil
	}

	var rawTTL string
	addr, rawTTL = addr[:slash], addr[slash+1:]
	if slash = strings.IndexByte(rawTTL, '/'); slash != -1 {
		addrrange, err = strconv.Atoi(rawTTL[slash+1:])
		if err != nil {
			return "", -1, 0, err
		}

		rawTTL = rawTTL[:slash]
	}

	ttl, err = strconv.Atoi(rawTTL)
	if err != nil {
		return "", -1, 0, err
	}

	if typ == IP6 {
		// Don't forget! IP6 doesn't have TTL.
		// Also, if we're at this point, so TTL is already presented. Then
		// just swap ttl with addrrange and set ttl to 0,
		// (!) ignoring possible error: both ttl and addrrange are presented
		addrrange = ttl
		ttl = -1
	}

	return addr, ttl, addrrange, nil
}
