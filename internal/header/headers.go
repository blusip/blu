package header

// headersPreAlloc is implemented, as we already know that headers will definitely be
// presented in every request. So we can a bit make the life of GC easier by avoiding
// extra map allocations (escapes, so this also positively affects cold performance)
const headersPreAlloc = 10

// Headers is a struct, that serves to encapsulate headers (SIP, SDP, etc.).
// This allows us to implement some optimizations easily, like swapping underlying
// implementation to more effective ones or optimizing by well-known headers
type Headers struct {
	headers map[string][]string
}

// NewHeaders returns a new instance of Headers with initialized underlying storage
func NewHeaders() Headers {
	return Headers{
		headers: make(map[string][]string, headersPreAlloc),
	}
}

// Get fetches the first value if presented, otherwise just an empty string
func (h Headers) Get(key string) (value string, found bool) {
	values, found := h.headers[key]
	if !found {
		return "", false
	}

	return values[0], true
}

// GetAll returns a complete slice of all the header values
func (h Headers) GetAll(key string) (values []string, found bool) {
	values, found = h.headers[key]
	return values, found
}

// Add appends a new value to the headers. In case key didn't exist before, a new entry
// will be created
func (h Headers) Add(key string, values ...string) {
	h.headers[key] = append(h.headers[key], values...)
}

// Set overrides the entry by provided values slice
func (h Headers) Set(key string, values ...string) {
	h.headers[key] = values
}

// Clear clears all the headers
func (h Headers) Clear() {
	// fun fact: this will be optimized into a single mapclear() call
	for k := range h.headers {
		delete(h.headers, k)
	}
}

// Unwrap returns the underlying implementation of Headers object.
func (h Headers) Unwrap() map[string][]string {
	return h.headers
}
