package ctxd

// FieldNames defines standard field names.
//
// Default names are aligned with https://www.elastic.co/guide/en/ecs/current/ecs-field-reference.html.
type FieldNames struct {
	Timestamp string `default:"@timestamp"`
	Message   string `default:"message"`

	// ClientIP is an IP address of the client (IPv4 or IPv6).
	ClientIP string `default:"client.ip"`

	HTTPMethod         string `default:"http.request.method"`
	HTTPResponseBytes  string `default:"http.response.bytes"`
	HTTPResponseStatus string `default:"http.response.status_code"`

	URL string `default:"url.original"`

	// UserAgentOriginal is an unparsed user_agent string.
	UserAgentOriginal string `default:"user_agent.original"`

	SpanID        string `default:"span.id"`
	TraceID       string `default:"trace.id"`
	TransactionID string `default:"transaction.id"`
}
