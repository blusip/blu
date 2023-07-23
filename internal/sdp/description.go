package sdp

type Description struct {
	Session Session
	Media   []Media
}

type Session struct {
	Protocol       string   // compulsory
	Originator     string   // compulsory
	Name           string   // compulsory
	Info           string   // optional
	URI            string   // optional
	Email          string   // optional
	Phone          string   // optional
	ConnectionData string   // optional
	BandwidthInfo  []string // optional

	// TODO: add time descriptions here

	TimeZoneAdjustments []string // optional
	EncryptionKey       string   // optional
	Attributes          []string // optional
}

type Media struct {
	Name           string   // compulsory
	Title          string   // optional
	ConnectionInfo string   // optional if included in session-level
	BandwidthInfo  []string // optional
	EncryptionKey  string   // optional
	Attributes     []string // optional
}
