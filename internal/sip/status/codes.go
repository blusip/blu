package status

type (
	Code   uint16
	Status string
)

//go:generate stringer -type=Code

// SIP status codes as registered RFC 3261
// See: https://www.rfc-editor.org/rfc/rfc3261
const (
	// Provisional 1xx codes
	Trying               Code = 100 // RFC 3261 21.1.1
	Ringing              Code = 180 // RFC 3261 21.1.2
	CallIsBeingForwarded Code = 181 // RFC 3261 21.1.3
	Queued               Code = 182 // RFC 3261 21.1.4
	SessionProgress      Code = 183 // RFC 3261 21.1.5

	// Successful 2xx codes
	OK       Code = 200 // RFC 3261 21.2.1
	Accepted Code = 201 // Not show in RFC but recently used as de-facto standart code

	// Redirection 3xx codes
	MultipleChoices    Code = 300 // RFC 3261 21.3.1
	MovedPermanently   Code = 301 // RFC 3261 21.3.2
	MovedTemporarily   Code = 302 // RFC 3261 21.3.3
	UseProxy           Code = 305 // RFC 3261 21.3.4
	AlternativeService Code = 380 // RFC 3261 21.3.5

	// Request Failure 4xx codes
	BadRequest                  Code = 400 // RFC 3261 21.4.1
	Unauthorized                Code = 401 // RFC 3261 21.4.2
	PaymentRequired             Code = 402 // RFC 3261 21.4.3
	Forbidden                   Code = 403 // RFC 3261 21.4.4
	NotFound                    Code = 404 // RFC 3261 21.4.5
	MethodNotAllowed            Code = 405 // RFC 3261 21.4.6
	NotAcceptable               Code = 406 // RFC 3261 21.4.7
	ProxyAuthenticationRequired Code = 407 // RFC 3261 21.4.8
	RequestTimeout              Code = 408 // RFC 3261 21.4.9
	Gone                        Code = 410 // RFC 3261 21.4.10
	RequestEntityTooLarge       Code = 413 // RFC 3261 21.4.11
	RequestURITooLong           Code = 414 // RFC 3261 21.4.12
	UnsupportedMediaType        Code = 415 // RFC 3261 21.4.13
	UnsupportedURIScheme        Code = 416 // RFC 3261 21.4.14
	BadExtension                Code = 420 // RFC 3261 21.4.15
	ExtensionRequired           Code = 421 // RFC 3261 21.4.16
	IntervalTooBrief            Code = 423 // RFC 3261 21.4.17
	TemporarilyUnavailable      Code = 480 // RFC 3261 21.4.18
	CallTransactionDoesNotExist Code = 481 // RFC 3261 21.4.19
	LoopDetected                Code = 482 // RFC 3261 21.4.20
	TooManyHops                 Code = 483 // RFC 3261 21.4.21
	AddressIncomplete           Code = 484 // RFC 3261 21.4.22
	Ambiguous                   Code = 485 // RFC 3261 21.4.23
	BusyHere                    Code = 486 // RFC 3261 21.4.24
	RequestTerminated           Code = 487 // RFC 3261 21.4.25
	NotAcceptableHere           Code = 488 // RFC 3261 21.4.26
	RequestPending              Code = 491 // RFC 3261 21.4.27
	Undecipherable              Code = 493 // RFC 3261 21.4.28

	// Server Failure 5xx codes
	ServerInternalError  Code = 500 // RFC 3261 21.5.1
	NotImplemented       Code = 501 // RFC 3261 21.5.2
	BadGateway           Code = 502 // RFC 3261 21.5.3
	ServiceUnavailable   Code = 503 // RFC 3261 21.5.4
	ServerTimeout        Code = 504 // RFC 3261 21.5.5
	VersionNotSupported  Code = 505 // RFC 3261 21.5.6
	MessageTooLarge      Code = 513 // RFC 3261 21.5.7
	URITooLong           Code = 520 // Custom 52x code
	URIDecoding          Code = 521 // Custom 52x code
	UnsupportedProtocol  Code = 522 // Custom 52x code
	TooManyHeaders       Code = 523 // Custom 52x code
	HeaderFieldsTooLarge Code = 524 // Custom 52x code
	MethodNotImplemented Code = 525 // Custom 52x code

	// Global Failures 6xx codes
	BusyEverywhere       Code = 600 // RFC 3261 21.6.1
	Decline              Code = 603 // RFC 3261 21.6.2
	DoesNotExistAnywhere Code = 604 // RFC 3261 21.6.3
	GlobalNotAcceptable  Code = 606 // RFC 3261 21.6.4
)

// Text returns a text for the SIP status code. It returns the empty
// string if the code is unknown.

func Text(code Code) Status {
	switch code {
	case Trying:
		return "Trying"

	case Ringing:
		return "Ringing"

	case CallIsBeingForwarded:
		return "Call Is Being Forwarded"

	case Queued:
		return "Queued"

	case OK:
		return "OK"

	case Accepted:
		return "Accepted"

	case MultipleChoices:
		return "Multiple Choices"

	case MovedPermanently:
		return "Moved Permanently"

	case MovedTemporarily:
		return "Moved Temporarily"

	case UseProxy:
		return "Use Proxy"

	case AlternativeService:
		return "Alternative Service"

	case BadRequest:
		return "Bad Request"

	case Unauthorized:
		return "Unauthorized"

	case PaymentRequired:
		return "Payment Required"

	case Forbidden:
		return "Forbidden"

	case MethodNotAllowed:
		return "Method Not Allowed"

	case NotAcceptable:
		return "Not Acceptable"

	case ProxyAuthenticationRequired:
		return "Proxy Authentication Required"

	case Gone:
		return "Gone"

	case RequestEntityTooLarge:
		return "Request Entity Too Large"

	case RequestURITooLong:
		return "Request URI Too Long"

	case UnsupportedMediaType:
		return "Unsupported Media Type"

	case UnsupportedURIScheme:
		return "Unsupported URI Scheme"

	case BadExtension:
		return "Bad Extension"

	case ExtensionRequired:
		return "Extension Required"

	case IntervalTooBrief:
		return "Interval Too Brief"

	case TemporarilyUnavailable:
		return "Temporarily Unavailable"

	case CallTransactionDoesNotExist:
		return "Call/Transaction Does Not Exist"

	case LoopDetected:
		return "Loop Detected"

	case TooManyHops:
		return "Too Many Hops"

	case AddressIncomplete:
		return "Address Incomplete"

	case Ambiguous:
		return "Ambiguous"

	case BusyHere:
		return "Busy Here"

	case RequestTerminated:
		return "Request Terminated"

	case NotAcceptableHere:
		return "Not Acceptable Here"

	case RequestPending:
		return "Request Pending"

	case Undecipherable:
		return "Undecipherable"

	case ServerInternalError:
		return "Server Internal Error"

	case NotImplemented:
		return "Not Implemented"

	case BadGateway:
		return "Bad Gateway"

	case ServiceUnavailable:
		return "Service Unavailable"

	case ServerTimeout:
		return "Server Timeout"

	case VersionNotSupported:
		return "Version Not Supported"

	case MessageTooLarge:
		return "Message Too Large"

	case URITooLong:
		return "URI Too Long"

	case URIDecoding:
		return "URI Decoding Failed"

	case TooManyHeaders:
		return "Too Many Headers"

	case HeaderFieldsTooLarge:
		return "Header Fields Too Large"

	case UnsupportedProtocol:
		return "Unsupported Protocol"

	case MethodNotImplemented:
		return "Method Not Implemented"

	case BusyEverywhere:
		return "Busy Everywhere"

	case Decline:
		return "Decline"

	case DoesNotExistAnywhere:
		return "Does Not Exist Anywhere"

	case GlobalNotAcceptable:
		return "Not Acceptable"

	default:
		return "Unknown Status Code"
	}
}

// CodeStatus returns a pre-defined line with code and status text (including
// terminating CRLF sequence) in case code is known to server, otherwise empty
// line is returned
func CodeStatus(code Code) string {
	switch code {
	case Trying:
		return "100 Trying\r\n"

	case Ringing:
		return "180 Ringing\r\n"

	case CallIsBeingForwarded:
		return "181 Call Is Being Forwarded\r\n"

	case Queued:
		return "182 Queued\r\n"

	case SessionProgress:
		return "183 Session Progress\r\n"

	case OK:
		return "200 OK\r\n"

	case Accepted:
		return "201 Accepted\r\n"

	case MultipleChoices:
		return "300 Multiple Choices\r\n"

	case MovedPermanently:
		return "301 Moved Permanently\r\n"

	case MovedTemporarily:
		return "302 Moved Temporarily\r\n"

	case UseProxy:
		return "305 Use Proxy\r\n"

	case AlternativeService:
		return "380 AlternativeS ervice\r\n"

	case BadRequest:
		return "400 Bad Request\r\n"

	case Unauthorized:
		return "401 Unauthorized\r\n"

	case PaymentRequired:
		return "402 Payment Required\r\n"

	case Forbidden:
		return "403 Forbidden\r\n"

	case NotFound:
		return "404 Not Found\r\n"

	case NotAcceptable:
		return "405 Not Acceptable\r\n"

	case MethodNotAllowed:
		return "406 Method Not Allowed\r\n"

	case ProxyAuthenticationRequired:
		return "407 Proxy Authentication Required\r\n"

	case RequestTimeout:
		return "408 Request Timeout\r\n"

	case Gone:
		return "410 Gone\r\n"

	case RequestEntityTooLarge:
		return "413 Request Entity Too Large\r\n"

	case RequestURITooLong:
		return "414 Request URI Too Long\r\n"

	case UnsupportedMediaType:
		return "415 Unsupported Media Type\r\n"

	case UnsupportedURIScheme:
		return "416 Unsupported URI Scheme\r\n"

	case BadExtension:
		return "420 Bad Extension\r\n"

	case ExtensionRequired:
		return "421 Extension Required\r\n"

	case IntervalTooBrief:
		return "423 Interval Too Brief\r\n"

	case TemporarilyUnavailable:
		return "480 Temporarily Unavailable\r\n"

	case CallTransactionDoesNotExist:
		return "481 Call/Transaction Does Not Exist\r\n"

	case LoopDetected:
		return "482 Loop Detected\r\n"

	case TooManyHops:
		return "483 Too Many Hops\r\n"

	case AddressIncomplete:
		return "484 Address Incomplete\r\n"

	case Ambiguous:
		return "485 Ambiguous\r\n"

	case BusyHere:
		return "486 Busy Here\r\n"

	case RequestTerminated:
		return "487 Request Terminated\r\n"

	case NotAcceptableHere:
		return "488 Not Acceptable Here\r\n"

	case RequestPending:
		return "491 Request Pending\r\n"

	case Undecipherable:
		return "493 Undecipherable\r\n"

	case ServerInternalError:
		return "500 Server Internal Error\r\n"

	case NotImplemented:
		return "501 Not Implemented\r\n"

	case BadGateway:
		return "502 Bad Gateway\r\n"

	case ServiceUnavailable:
		return "503 Service Unavailable\r\n"

	case ServerTimeout:
		return "504 Server Timeout\r\n"

	case VersionNotSupported:
		return "505 Version Not Supported\r\n"

	case MessageTooLarge:
		return "513 Message Too Large\r\n"

	case URITooLong:
		return "520 URI Too Long\r\n"

	case URIDecoding:
		return "521 URI Decoding Failed\r\n"

	case UnsupportedProtocol:
		return "522 Unsupported Protocol\r\n"

	case TooManyHeaders:
		return "523 Too Many Headers\r\n"

	case HeaderFieldsTooLarge:
		return "524 Header Fields Too Large\r\n"

	case MethodNotImplemented:
		return "525 Method Not Implemented\r\n"

	case BusyEverywhere:
		return "600 Busy Everywhere\r\n"

	case Decline:
		return "603 Decline\r\n"

	case DoesNotExistAnywhere:
		return "604 Does Not Exist Anywhere\r\n"

	case GlobalNotAcceptable:
		return "606 Not Acceptable\r\n"

	default:
		return ""
	}
}
