package parser

import (
	"github.com/gokiki/sip-server/internal/sip"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParser(t *testing.T) {
	t.Run("request from RFC", func(t *testing.T) {
		data := "" +
			"INVITE sip:bob@biloxi.com SIP/2.0\r\n" +
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
		p := NewParser(request)
		done, err := p.Parse([]byte(data))
		require.NoError(t, err)
		require.True(t, done, "given the whole request at once, parser is expected to be done")

		require.Equal(t, "INVITE", request.Method)
	})
}
