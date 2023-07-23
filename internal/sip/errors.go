package sip

type Error struct {
	Message string
	Code    Code
}

func NewError(code Code, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}

func (e Error) Error() string {
	return e.Message
}

var (
	ErrBadRequest                  = NewError(BadRequest, "bad request")
	ErrUnauthorized                = NewError(Unauthorized, "unauthorized")
	ErrForbidden                   = NewError(Forbidden, "forbidden")
	ErrMethodNotAllowed            = NewError(MethodNotAllowed, "method not allowed")
	ErrNotAcceptable               = NewError(NotAcceptable, "not acceptable")
	ErrProxyAuthenticationRequired = NewError(ProxyAuthenticationRequired, "proxy authentication required")
	ErrGone                        = NewError(Gone, "gone")
	ErrRequestEntityTooLarge       = NewError(RequestEntityTooLarge, "request entity too large")
	ErrRequestURITooLong           = NewError(RequestURITooLong, "request URI too long")
	ErrUnsupportedMediaType        = NewError(UnsupportedMediaType, "unsupported media type")
	ErrUnsupportedURIScheme        = NewError(UnsupportedURIScheme, "bad unsupported URI scheme")
	ErrExtensionRequired           = NewError(ExtensionRequired, "extension required")
	ErrIntervalTooBrief            = NewError(IntervalTooBrief, "interval too brief")
	ErrTemporarilyUnavailable      = NewError(TemporarilyUnavailable, "temporarily unavailable")
	ErrCallTransactionDoesNotExist = NewError(CallTransactionDoesNotExist, "call/transaction does not exist")
	ErrLoopDetected                = NewError(LoopDetected, "loop detected")
	ErrTooManyHops                 = NewError(TooManyHops, "too many hops")
	ErrAddressIncomplete           = NewError(AddressIncomplete, "address incomplete")
	ErrAmbiguous                   = NewError(Ambiguous, "ambiguous")
	ErrBusyHere                    = NewError(BusyHere, "busy here")
	ErrRequestTerminated           = NewError(RequestTerminated, "request terminated")
	ErrNotAcceptableHere           = NewError(NotAcceptableHere, "not acceptable here")
	ErrRequestPending              = NewError(RequestPending, "request pending")
	ErrUndecipherable              = NewError(Undecipherable, "undecipherable")
	ErrServerInternalError         = NewError(ServerInternalError, "server internal error")
	ErrNotImplemented              = NewError(NotImplemented, "not implemented")
	ErrBadGateway                  = NewError(BadGateway, "bad gateway")
	ErrServiceUnavailable          = NewError(ServiceUnavailable, "service unavailable")
	ErrServerTimeout               = NewError(ServerTimeout, "server timeout")
	ErrVersionNotSupported         = NewError(VersionNotSupported, "version not supported")
	ErrMessageTooLarge             = NewError(MessageTooLarge, "message too large")
	ErrMethodNotImplemented        = NewError(MethodNotImplemented, "method not implemented")
	ErrURITooLong                  = NewError(URITooLong, "URI too long")
	ErrURIDecoding                 = NewError(URIDecoding, "URI decoding failed")
	ErrUnsupportedProtocol         = NewError(UnsupportedProtocol, "unsupported protocol")
	ErrTooManyHeaders              = NewError(TooManyHeaders, "too many headers")
	ErrHeaderFieldsTooLarge        = NewError(HeaderFieldsTooLarge, "header fields too large")
	ErrBusyEverywhere              = NewError(BusyEverywhere, "busy everywhere")
	ErrDecline                     = NewError(Decline, "decline")
	ErrDoesNotExistAnywhere        = NewError(DoesNotExistAnywhere, "does not exist anywhere")
	ErrGlobalNotAcceptable         = NewError(GlobalNotAcceptable, "not acceptable")
)
