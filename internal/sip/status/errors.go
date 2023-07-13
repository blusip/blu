package status

type SIPError struct {
	Message string
	Code    Code
}

func newErr(code Code, message string) SIPError {
	return SIPError{
		Code:    code,
		Message: message,
	}
}

func (h SIPError) Error() string {
	return h.Message
}

var (
	ErrBadRequest                  = newErr(BadRequest, "bad request")
	ErrUnauthorized                = newErr(Unauthorized, "unauthorized")
	ErrForbidden                   = newErr(Forbidden, "forbidden")
	ErrMethodNotAllowed            = newErr(MethodNotAllowed, "method not allowed")
	ErrNotAcceptable               = newErr(NotAcceptable, "not acceptable")
	ErrProxyAuthenticationRequired = newErr(ProxyAuthenticationRequired, "proxy authentication required")
	ErrGone                        = newErr(Gone, "gone")
	ErrRequestEntityTooLarge       = newErr(RequestEntityTooLarge, "request entity too large")
	ErrRequestURITooLong           = newErr(RequestURITooLong, "request URI too long")
	ErrUnsupportedMediaType        = newErr(UnsupportedMediaType, "unsupported media type")
	ErrUnsupportedURIScheme        = newErr(UnsupportedURIScheme, "bad unsupported URI scheme")
	ErrExtensionRequired           = newErr(ExtensionRequired, "extension required")
	ErrIntervalTooBrief            = newErr(IntervalTooBrief, "interval too brief")
	ErrTemporarilyUnavailable      = newErr(TemporarilyUnavailable, "temporarily unavailable")
	ErrCallTransactionDoesNotExist = newErr(CallTransactionDoesNotExist, "call/transaction does not exist")
	ErrLoopDetected                = newErr(LoopDetected, "loop detected")
	ErrTooManyHops                 = newErr(TooManyHops, "too many hops")
	ErrAddressIncomplete           = newErr(AddressIncomplete, "address incomplete")
	ErrAmbiguous                   = newErr(Ambiguous, "ambiguous")
	ErrBusyHere                    = newErr(BusyHere, "busy here")
	ErrRequestTerminated           = newErr(RequestTerminated, "request terminated")
	ErrNotAcceptableHere           = newErr(NotAcceptableHere, "not acceptable here")
	ErrRequestPending              = newErr(RequestPending, "request pending")
	ErrUndecipherable              = newErr(Undecipherable, "undecipherable")
	ErrServerInternalError         = newErr(ServerInternalError, "server internal error")
	ErrNotImplemented              = newErr(NotImplemented, "not implemented")
	ErrBadGateway                  = newErr(BadGateway, "bad gateway")
	ErrServiceUnavailable          = newErr(ServiceUnavailable, "service unavailable")
	ErrServerTimeout               = newErr(ServerTimeout, "server timeout")
	ErrVersionNotSupported         = newErr(VersionNotSupported, "version not supported")
	ErrMessageTooLarge             = newErr(MessageTooLarge, "message too large")
	ErrURITooLong                  = newErr(URITooLong, "URI too long")
	ErrURIDecoding                 = newErr(URIDecoding, "URI decoding failed")
	ErrUnsupportedProtocol         = newErr(UnsupportedProtocol, "unsupported protocol")
	ErrTooManyHeaders              = newErr(TooManyHeaders, "too many headers")
	ErrHeaderFieldsTooLarge        = newErr(HeaderFieldsTooLarge, "header fields too large")
	ErrBusyEverywhere              = newErr(BusyEverywhere, "busy everywhere")
	ErrDecline                     = newErr(Decline, "decline")
	ErrDoesNotExistAnywhere        = newErr(DoesNotExistAnywhere, "does not exist anywhere")
	ErrGlobalNotAcceptable         = newErr(GlobalNotAcceptable, "not acceptable")
)
