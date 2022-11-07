package middleware

import (
	"fmt"
	"net/http"
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

	status := c.StatusCode()
	if err != nil {
		status = http.StatusInternalServerError
		fileds["error"] = err.Error()
	}

	// if err != nil {
	// 	switch e := err.(type) {
	// 	case *fiber.Error:
	// 		status = e.Code
	// 	default: // case error
	// 		status = fiber.StatusInternalServerError
	// 	}

	// 	fileds["error"] = err.Error()
	// 	log.Error(err.Error())
	// }

	fileds["status"] = c.StatusCode()
	fileds["latency"] = time.Since(start)

	msg := fmt.Sprintf("%d - %s %s", status, c.Method(), c.Path())

	logger.New(logger.ToFields(fileds)...).Info(msg)

	return err
}
