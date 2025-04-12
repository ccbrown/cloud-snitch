package api

type Config struct {
	// The number of proxies positioned in front of the API. This is used to interpret
	// X-Forwarded-For headers.
	ProxyCount int

	// If set, X-Forwarded-For headers are ignored unless there is also a "Proxy-Secret" header with this
	// value.
	ProxySecret string
}
