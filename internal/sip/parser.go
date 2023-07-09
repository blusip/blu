package sip

type Parser interface {
	Parse(data []byte) (finished bool, err error)
	Finish()
}
