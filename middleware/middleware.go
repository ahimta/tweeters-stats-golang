package middleware

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Ahimta/tweeters-stats-golang/config"
)

const (
	// https://httpd.apache.org/docs/2.2/logs.html#combined + execution time.
	apacheFormatPattern = "%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\" %.3f\n"
	xForwardedFor       = "X-Forwarded-For"
)

// Apply applies generic middlware to an http.Handler
// currently it includes:
// * logging middleware
// * error-handling middleware
// * security middleware
// * CORS middleware
// * CSRF middleware
func Apply(
	handler http.Handler,
	writer io.Writer,
	c *config.Config,
) http.Handler {

	host := c.Host
	origin := fmt.Sprintf("%s://%s", c.Protocol, host)
	referrerPrefix := fmt.Sprintf("%s/", origin)

	return http.HandlerFunc(func(w0 http.ResponseWriter, r *http.Request) {
		// logging middleware
		startTime := time.Now()
		w := &responseWriter{w0, 200, 0}

		// error-handling middleware
		defer func() {
			r := recover()
			var err error

			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}

				log.Println("panic", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
			}
		}()

		// logging middleware
		clientIP := r.RemoteAddr
		if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
			clientIP = clientIP[:colon]
		}

		if s := r.Header.Get(xForwardedFor); s != "" {
			clientIP = s
		}

		referer := r.Referer()
		if referer == "" {
			referer = "-"
		}

		userAgent := r.UserAgent()
		if userAgent == "" {
			userAgent = "-"
		}

		defer func() {
			finishTime := time.Now()
			time := finishTime.UTC()
			elapsedTime := finishTime.Sub(startTime)
			timeFormatted := time.Format("02/Jan/2006 03:04:05")

			status := w.status
			responseBytes := w.responseBytes

			fmt.Fprintf(
				writer,
				apacheFormatPattern,
				clientIP,
				timeFormatted,
				r.Method,
				r.URL,
				r.Proto,
				status,
				responseBytes,
				referer,
				userAgent,
				elapsedTime.Seconds(),
			)
		}()

		// Security middleware
		w.Header().Set("Content-Security-Policy", "default-src 'self' data: maxcdn.bootstrapcdn.com; style-src 'unsafe-inline' maxcdn.bootstrapcdn.com; script-src 'unsafe-inline'")
		w.Header().Set("Referrer-Header", "same-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=5184000")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		w.Header().Set("X-Download-Options", "noopen")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// CORS middleare
		if c.CorsDomain != "" {
			w.Header().Set("Access-Control-Allow-Origin", c.CorsDomain)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		}

		// CSRF middleware
		if path := r.URL.Path; !(path == "/" ||
			path == "/login/twitter" ||
			path == "/oauth/twitter/callback") {

			if r.Header.Get("X-Requested-With") != "XMLHttpRequest" ||
				r.Host != host {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if !(r.Header.Get("Origin") == origin ||
				r.Header.Get("Referer") == origin ||
				strings.HasPrefix(r.Header.Get("Referer"), referrerPrefix)) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		// related to CORS middleware
		if r.Method != http.MethodOptions {
			handler.ServeHTTP(w, r)
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	status        int
	responseBytes int64
}

func (w *responseWriter) Write(p []byte) (int, error) {
	written, err := w.ResponseWriter.Write(p)
	w.responseBytes += int64(written)
	return written, err
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
