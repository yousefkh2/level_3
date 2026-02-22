package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yousefkh2/level_3/week_4/api/app"
	"go.uber.org/zap"
)

func Observability(application *app.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			start := time.Now()

			application.Metrics.RequestsInFlight.Inc()
			defer application.Metrics.RequestsInFlight.Dec()

			err := next(c)

			duration := time.Since(start)
			status := c.Response().Status
			method := req.Method

			path := c.Path()
			if path == "" {
				path = req.URL.Path
			}

			statusStr := fmt.Sprintf("%d", status)

			application.Metrics.RequestsTotal.WithLabelValues(method, path, statusStr).Inc()

			application.Metrics.RequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())

			fields := []zap.Field{
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", status),
				zap.Duration("duration", duration),
				zap.String("remote_ip", c.RealIP()),
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
				application.Logger.Error("request failed", fields...)
			} else if status >= 500 {
				application.Logger.Error("request completed with server error", fields...)
			} else if status >= 400 {
				application.Logger.Warn("request completed with client error", fields...)
			} else {
				application.Logger.Info("request completed", fields...)
			}

			return err
		}
	}
}
