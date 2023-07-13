package parser

import (
	"github.com/gokiki/sip-server/internal/sip"
	"github.com/gokiki/sip-server/settings"
	"github.com/indigo-web/utils/arena"
	"github.com/indigo-web/utils/pool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func newParser(request *sip.Request) *Parser {
	keyArena := *arena.NewArena[byte](0, 65535)
	valArena := *arena.NewArena[byte](0, 65535)
	requestLineArena := *arena.NewArena[byte](0, 65535)
	valuesPool := *pool.NewObjectPool[[]string](10)

	return NewParser(request, keyArena, valArena, requestLineArena, valuesPool, settings.Default())
}

func TestParser(t *testing.T) {
	t.Run("default request", func(t *testing.T) {
		data := "" +
			"INVITE sip:bob%20smith:fancy%20password@biloxi.com:80;par+am=value;par%20am2=val%20ue2 SIP/2.0\r\n" +
			"Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds\r\n" +
			"Max-Forwards: 70\r\n" +
			"To: Bob <sip:bob@biloxi.com>\r\n" +
			"From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n" +
			"Call-ID: a84b4c76e66710@pc33.atlanta.com\r\n" +
			"CSeq: 314159 INVITE\r\n" +
			"Contact: <sip:alice@pc33.atlanta.com>\r\n" +
			"Content-Type: application/sdp\r\n" +
			"Content-Length: 13\r\n\r\n" +
			"some SDP here"

		request := sip.NewRequest()
		p := newParser(request)
		done, err := p.Parse([]byte(data))
		require.NoError(t, err)
		require.True(t, done, "given the whole request at once, parser is expected to be done")
		require.Equal(t, "INVITE", request.Method)
		require.Equal(t, "sip", request.URI.Scheme)
		require.Equal(t, "bob smith", request.URI.User)
		require.Equal(t, "fancy password", request.URI.Password)
		require.Equal(t, "biloxi.com", request.URI.Host)
		require.Equal(t, 80, request.URI.Port)
		value, found := request.URI.Params.Get("par am")
		require.Truef(t, found, "wanted \"par am\" parameter")
		require.Equal(t, "value", value)
		value, found = request.URI.Params.Get("par am2")
		require.Truef(t, found, "wanted \"par am2\" parameter")
		require.Equal(t, "val ue2", value)
		require.Equal(t, "SIP", request.Proto.Scheme())
		require.Equal(t, "2.0", request.Proto.Version())
		require.Equal(t, 13, request.ContentLength)

		for key, want := range map[string]string{
			"Via":          "SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds",
			"Max-Forwards": "70",
			"To":           "Bob <sip:bob@biloxi.com>",
			"From":         "Alice <sip:alice@atlanta.com>;tag=1928301774",
			"Call-ID":      "a84b4c76e66710@pc33.atlanta.com",
			"CSeq":         "314159 INVITE",
			"Contact":      "<sip:alice@pc33.atlanta.com>",
			"Content-Type": "application/sdp",
		} {
			value, found := request.Headers.Get(key)
			if !assert.Truef(t, found, "header not found: %s", key) {
				continue
			}

			assert.Equalf(t, want, value, "header's values doesn't match: %s", key)
		}
	})
}
