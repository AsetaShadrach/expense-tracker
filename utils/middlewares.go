package utils

import (
	"bytes"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
)

func InstrumentRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path)
		defer span.End()

		w.Header().Set("trace-id", span.SpanContext().TraceID().String())

		attrs := HeaderToAttributes(r.Header)
		span.SetAttributes(attrs...)

		bodyBytes, _ := io.ReadAll(r.Body)

		// Close the original body
		r.Body.Close()

		// Instrument the body content (e.g., as a string)
		span.SetAttributes(attribute.String("http.request.body", string(bodyBytes)))

		// Create a new ReadCloser from the buffered bytes and assign it back to the request
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		next.ServeHTTP(w, r)
	})
}
