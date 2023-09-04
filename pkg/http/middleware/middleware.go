package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/s-gurman/user-segmentation/pkg/logger"
)

func PanicRecovery(next http.Handler, l logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println()
		l.Info("panic recovery middleware")

		defer func() {
			if err := recover(); err != nil {
				l.Errorf("panic recovered with err: %s", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func AccessLog(next http.Handler, l logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("access log middleware")

		start := time.Now()
		next.ServeHTTP(w, r)

		l.Infow("served",
			"method", r.Method,
			"uri", r.RequestURI,
			"from", r.RemoteAddr,
			"elapsed", fmt.Sprint(time.Since(start)),
		)
	})
}
