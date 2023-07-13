package settings

import (
	"math"
)

func Default() Settings {
	return Settings{
		RequestLine: RequestLine{
			MaxLength:      65535,
			BufferPreAlloc: 1024,
		},
		Headers: Headers{
			MaxNumber:      200,
			MaxKeyLength:   32768,
			MaxValueLength: 65535 * 2,
		},
		Body: Body{
			MaxLength:      math.MaxInt32,
			BufferPreAlloc: 1024,
		},
	}
}
