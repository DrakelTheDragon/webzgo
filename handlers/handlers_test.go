package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubMiddleware struct {
	logger *log.Logger
}

const stubLogPrefix = "Received Request:"

const stubLogFormat = `
  - Method: %s
  - RemoteAddr: %s
  - Host: %s
  - URL: %s
`

func (s *StubMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s.logger.Printf(stubLogFormat, r.Method, r.RemoteAddr, r.Host, r.URL)

			next.ServeHTTP(w, r)
		})
}

func TestWrap(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	response := httptest.NewRecorder()

	var buf bytes.Buffer

	mw := &StubMiddleware{logger: log.New(&buf, stubLogPrefix, 0)}

	handler := mw.Wrap(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(w, "Hello, World!")
		},
	))

	handler.ServeHTTP(response, request)

	assertRequestLog(t, request, &buf)

	got := response.Body.String()
	want := "Hello, World!"

	assertResponseBody(t, got, want)
}

func assertResponseBody(t testing.TB, got, want string) {
	if got != want {
		t.Errorf("\nhave: %q\nwant: %q", got, want)
	}
}

func assertRequestLog(t testing.TB, r *http.Request, buf *bytes.Buffer) {
	t.Helper()

	msg := fmt.Sprintf(stubLogFormat, r.Method, r.RemoteAddr, r.Host, r.URL)

	got := buf.String()
	want := stubLogPrefix + msg

	if got != want {
		t.Errorf("\nhave: %q\nwant: %q", got, want)
	}
}
