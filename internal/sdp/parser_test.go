package sdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Run("rfc sample", func(t *testing.T) {
		rfcSample := []byte("v=0\r\n" +
			"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
			"s=SDP Seminar\r\n" +
			"i=A Seminar on the session description protocol\r\n" +
			"u=http://www.example.com/seminars/sdp.pdf\r\n" +
			"e=j.doe@example.com (Jane Doe)\r\n" +
			"c=IN IP4 224.2.17.12/127\r\n" +
			"b=CT:128\r\n" +
			//"t=2873397496 2873404696\r\n" +
			"a=recvonly\r\n" +
			"m=audio 49170 RTP/AVP 0\r\n" +
			"m=video 51372 RTP/AVP 99\r\n" +
			"a=rtpmap:99 h263-1998/90000\r\n")
		parser := NewParser()
		desc, err := parser.Parse(rfcSample)
		require.NoError(t, err)
		require.Equal(t, "0", desc.Session.Protocol)
		require.Equal(t, "jdoe", desc.Session.Originator.Username)
		require.Equal(t, "2890844526", desc.Session.Originator.SessId)
		require.Equal(t, "2890842807", desc.Session.Originator.SessVersion)
		require.Equal(t, IN, desc.Session.Originator.NetType)
		require.Equal(t, IP4, desc.Session.Originator.AddrType)
		require.Equal(t, "10.47.16.5", desc.Session.Originator.UnicastAddress)
		require.Equal(t, "SDP Seminar", desc.Session.Name)
		require.Equal(t, "A Seminar on the session description protocol", desc.Session.Info)
		require.Equal(t, "http://www.example.com/seminars/sdp.pdf", desc.Session.URI)
		require.Equal(t, "j.doe@example.com (Jane Doe)", desc.Session.Email)
		require.Equal(t, 1, len(desc.Session.ConnectionInfo), "wanted single connection info")
		require.Equal(t, IN, desc.Session.ConnectionInfo[0].NetType)
		require.Equal(t, IP4, desc.Session.ConnectionInfo[0].AddrType)
		require.Equal(t, 127, desc.Session.ConnectionInfo[0].TTL)
		require.Equal(t, 0, desc.Session.ConnectionInfo[0].AddrRange)
		require.Equal(t, "224.2.17.12", desc.Session.ConnectionInfo[0].Address)
		require.Equal(t, 1, len(desc.Session.BandwidthInfo), "wanted single bandwidth info")
		require.Equal(t, CT, desc.Session.BandwidthInfo[0].Type)
		require.Equal(t, 128, desc.Session.BandwidthInfo[0].Value)
		require.Equal(t, []string{"recvonly"}, desc.Session.Attributes)
		require.Equal(t, 2, len(desc.Media), "must be exactly 2 media blocks")
		media := desc.Media[0]
		require.Equal(t, "audio 49170 RTP/AVP 0", media.Name)
		media = desc.Media[1]
		require.Equal(t, "video 51372 RTP/AVP 99", media.Name)
		require.Equal(t, []string{"rtpmap:99 h263-1998/90000"}, media.Attributes)
	})
}

func TestParseAddr(t *testing.T) {
	baseAddr := "127.0.0.1"

	t.Run("IP4 address", func(t *testing.T) {
		sample := baseAddr
		addr, ttl, addrrange, err := parseAddress(sample, IP4)
		require.NoError(t, err)
		require.Equal(t, -1, ttl)
		require.Zero(t, addrrange)
		require.Equal(t, sample, addr)
	})

	t.Run("IP4 address and TTL", func(t *testing.T) {
		sample := baseAddr + "/127"
		addr, ttl, addrrange, err := parseAddress(sample, IP4)
		require.NoError(t, err)
		require.Equal(t, 127, ttl)
		require.Zero(t, addrrange)
		require.Equal(t, baseAddr, addr)
	})

	t.Run("IP4 address, TTL and range", func(t *testing.T) {
		sample := baseAddr + "/127/4"
		addr, ttl, addrrange, err := parseAddress(sample, IP4)
		require.NoError(t, err)
		require.Equal(t, 127, ttl)
		require.Equal(t, 4, addrrange)
		require.Equal(t, baseAddr, addr)
	})

	t.Run("IP6 address and TTL", func(t *testing.T) {
		sample := baseAddr + "/127"
		addr, ttl, addrrange, err := parseAddress(sample, IP6)
		require.NoError(t, err)
		require.Equal(t, -1, ttl)
		require.Equal(t, 127, addrrange)
		require.Equal(t, baseAddr, addr)
	})
}
