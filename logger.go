package ginlogrus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const timeFormat = "02/Jan/2006:15:04:05 -0700"

// Logger is the logrus logger handler
func Logger(logger logrus.FieldLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := logger.WithFields(logrus.Fields{
			"statusCode":     statusCode,
			"duration":       stop.Nanoseconds(), // in nanoseconds
			"durationPretty": stop.String(),
			"clientIP":       clientIP,
			"method":         c.Request.Method,
			"path":           path,
			"referer":        referer,
			"dataLength":     dataLength,
			"userAgent":      clientUserAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf(
				"%s - [%s] \"%s %s\" %d %d \"%s\" \"%s\" (%s)",
				clientIP,
				time.Now().Format(timeFormat),
				c.Request.Method,
				path,
				statusCode,
				dataLength,
				referer,
				clientUserAgent,
				stop.String(),
			)
			if statusCode >= http.StatusInternalServerError {
				entry.Error(msg)
			} else if statusCode >= http.StatusBadRequest {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
