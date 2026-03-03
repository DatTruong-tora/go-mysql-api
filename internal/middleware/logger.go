package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		logrus.WithFields(logrus.Fields{
			"method":  r.Method,
			"path":    r.URL.Path,
			"latency": time.Since(start).String(),
			"ip":      r.RemoteAddr,
		}).Info("The end of HTTP Request")
	})
}
