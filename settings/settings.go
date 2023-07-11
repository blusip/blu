package settings

type (
	RequestLine struct {
		MaxLength       int
		BufferPreAlloc  int
		MaxMethodLength int
	}

	Headers struct {
		MaxNumber      int
		MaxKeyLength   int
		MaxValueLength int
	}

	Body struct {
		MaxLength      int
		BufferPreAlloc int
	}
)

type Settings struct {
	RequestLine RequestLine
	Headers     Headers
	Body        Body
}
