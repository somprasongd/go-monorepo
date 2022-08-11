package common

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/somprasongd/go-monorepo/common/logger"
)

func LoggerMiddleware(c HContext) error {
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

type TokenUser struct {
	UserId   string `json:"user_id"`
	Identity string `json:"identity"` // email or username
	Role     string `json:"role"`
}

func EncodeUserMiddleware(c HContext) error {
	log := c.Locals("log").(logger.Interface)

	idToken := c.Get("X-Id-Token")

	if idToken != "" {
		return ResponseError(c, ErrNotAllowIdToken)
	}

	cliams := c.Locals("claims").(jwt.MapClaims)

	tu := TokenUser{
		UserId:   cliams["user_id"].(string),
		Identity: cliams["email"].(string),
		Role:     cliams["role"].(string),
	}

	jsonStr, err := json.Marshal(tu)
	if err != nil {
		log.Error(err.Error())
		return ResponseError(c, ErrInvalidIdToken)
	}

	idToken = Base64Encode(string(jsonStr))

	c.Set("X-Id-Token", idToken)

	return c.Next()
}

func DecodeUserMiddleware(c HContext) error {
	log := c.Locals("log").(logger.Interface)

	idToken := c.Get("X-Id-Token")

	if idToken == "" {
		return ResponseError(c, ErrNoIdToken)
	}

	jsonStr, ok := Base64Decode(idToken)
	if !ok {
		return ResponseError(c, ErrInvalidIdToken)
	}

	fmt.Println(jsonStr)

	tu := TokenUser{}
	err := json.Unmarshal([]byte(jsonStr), &tu)
	if err != nil {
		log.Error(err.Error())
		return ResponseError(c, ErrInvalidIdToken)
	}

	c.Locals("user", tu)

	return c.Next()
}
