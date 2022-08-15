package middleware

import (
	"os"
	"time"

	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/common/logger"
)

func LoggerMiddleware(c common.HContext) error {
	start := time.Now()

	appName := os.Getenv("APP_NAME")

	if len(appName) == 0 {
		appName = "goapi"
	}

	fileds := map[string]interface{}{
		"app":       appName,
		"domain":    c.Domain(),
		"requestId": c.RequestId(),
		"userAgent": c.Get("User-Agent"),
		"ip":        c.ClientIP(),
		"method":    c.Method(),
		"traceId":   c.Get("X-B3-Traceid"),
		"spanId":    c.Get("X-B3-Spanid"),
		"uri":       c.Path(),
	}

	log := logger.New(logger.ToFields(fileds)...)

	c.Locals("log", log)

	err := c.Next()

	// "status - method path (latency)"
	// msg := fmt.Sprintf("%v - %v %v (%v)", c.StatusCode(), c.Method(), c.Path(), time.Since(start))

	fileds["status"] = c.StatusCode()
	fileds["latency"] = time.Since(start)

	logger.New(logger.ToFields(fileds)...).Info("")

	return err
}
