package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime/debug"

	"github.com/AsetaShadrach/expense-tracker/schemas"
	"go.opentelemetry.io/otel/attribute"
)

func InstrumentRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path)
		w.Header().Set("trace-id", span.SpanContext().TraceID().String())

		defer span.End()

		attrs := HeaderToAttributes("http.request.headers", r.Header)
		span.SetAttributes(attrs...)

		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body.Close()

		// Instrument the body content (e.g., as a string)
		span.SetAttributes(attribute.String("http.request.body", string(bodyBytes)))

		// Create a new ReadCloser from the buffered bytes and assign it back to the request
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Create a new recoder for the response
		newRw := httptest.NewRecorder()
		newRw.Header().Write(w)

		next.ServeHTTP(newRw, r)

		respBytes, _ := io.ReadAll(newRw.Result().Body)
		span.SetAttributes(attribute.String("http.response.body", string(respBytes)))

		responseAttrs := HeaderToAttributes("http.response.headers", newRw.Header())
		span.SetAttributes(responseAttrs...)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(newRw.Result().StatusCode)
		w.Write(respBytes)
	})
}

func ErrorResolver(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				GeneralLogger.Error(fmt.Sprintf("Exception occured : % v", string(debug.Stack())))
				_, span := tracer.Start(r.Context(), "exceptionOccured")
				defer span.End()

				errorTrace := make(map[string]interface{})
				errorTrace["exception.stack.trace"] = string(debug.Stack())

				span.SetAttributes(MapToAttributes(errorTrace)...)
				w.Header().Set("trace-id", span.SpanContext().TraceID().String())

				tracedErrors := schemas.ErrorList{
					ResponseCode: "EX001",
					Message:      "An error occured",
					Errors:       []string{"Internal server Error"},
				}
				byts, _ := json.Marshal(tracedErrors)
				http.Error(w, string(byts), http.StatusInternalServerError)
				// w.Write(byts)
				// w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
