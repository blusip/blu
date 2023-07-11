package sip

type Method string

const (
	INVITE   Method = "INVITE"
	REGISTER Method = "REGISTER"
	OPTIONS  Method = "OPTIONS"
	ACK      Method = "ACK"
	CANCEL   Method = "CANCEL"
	BYE      Method = "BYE"
)
