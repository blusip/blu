package sdp

type Description struct {
	Session Session
	Media   []Media
}

type Session struct {
	Protocol       string           // compulsory
	Originator     Origin           // compulsory
	Name           string           // compulsory
	Info           string           // optional
	URI            string           // optional
	Email          string           // optional
	Phone          string           // optional
	ConnectionInfo []ConnectionInfo // optional
	BandwidthInfo  []Bandwidth      // optional

	// TODO: add time descriptions here

	TimeZoneAdjustments []string      // optional
	EncryptionKey       EncryptionKey // optional
	Attributes          []Attribute   // optional
}

type Media struct {
	Name           string           // compulsory
	Title          string           // optional
	ConnectionInfo []ConnectionInfo // optional if included in session-level
	BandwidthInfo  []string         // optional
	EncryptionKey  string           // optional
	Attributes     []string         // optional
}
