package sip

import "github.com/gokiki/sip-server/internal/header"

type URI struct {
	Scheme   string
	User     string
	Password string
	Host     string
	Port     int
	// TODO: give more generic name to header.Headers, as, how we can see here, it is used
	//  not only for headers, but as general key-value structure
	Params header.Headers
}

type Protocol string

func (p Protocol) Scheme() string {
	for i := range p {
		if p[i] == '/' {
			return string(p[:i])
		}
	}

	return string(p)
}

func (p Protocol) Version() string {
	for i := range p {
		if p[i] == '/' {
			return string(p[i+1:])
		}
	}

	return ""
}

type Request struct {
	Method        string
	URI           URI
	Proto         Protocol
	Headers       header.Headers
	ContentLength int
	Body          []byte
}

func NewRequest() *Request {
	return &Request{
		Headers: header.NewHeaders(),
		URI: URI{
			Params: header.NewHeaders(),
		},
	}
}

func (r Request) HasBody() bool {
	return r.ContentLength > 0
}
