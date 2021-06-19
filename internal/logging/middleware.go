package logging

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// Middleware logs the incoming HTTP request & its duration.
func Middleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.WithFields(
						log.Fields{
							"status": http.StatusInternalServerError,
							"method": r.Method,
							"path":   r.URL.Path,
							"query":  r.URL.RawQuery,
						},
					).Error(fmt.Sprintf("Error: %v\n%s", err, debug.Stack()))
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.WithFields(
				log.Fields{
					"status":   wrapped.status,
					"method":   r.Method,
					"path":     r.URL.EscapedPath(),
					"query":    r.URL.RawQuery,
					"duration": time.Since(start),
				},
			)
		}

		return http.HandlerFunc(fn)
	}
}
