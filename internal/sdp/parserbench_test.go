package sdp

import "testing"

func BenchmarkParser(b *testing.B) {
	rfcSample := []byte("v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"s=SDP Seminar\r\n" +
		"i=A Seminar on the session description protocol\r\n" +
		"u=http://www.example.com/seminars/sdp.pdf\r\n" +
		"e=j.doe@example.com (Jane Doe)\r\n" +
		"c=IN IP4 224.2.17.12/127\r\n" +
		//"t=2873397496 2873404696\r\n" +
		"a=recvonly\r\n" +
		"m=audio 49170 RTP/AVP 0\r\n" +
		"m=video 51372 RTP/AVP 99\r\n" +
		"a=rtpmap:99 h263-1998/90000\r\n")

	b.Run("rfc sample", func(b *testing.B) {
		parser := NewParser()
		b.SetBytes(int64(len(rfcSample)))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = parser.Parse(rfcSample)
		}
	})
}
