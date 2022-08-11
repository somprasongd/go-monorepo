package common

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/somprasongd/go-monorepo/common/logger"
)

func LoggerMiddleware(c HContext) error {
	start := time.Now()

	fileds := map[string]interface{}{}
	fileds["ip"] = c.IP()
	fileds["port"] = c.Port()
	fileds["requestid"] = c.RequestId()

	log := logger.NewWithFields(fileds)

	c.Locals("log", log)

	err := c.Next()

	// "status - method path (duration)"
	msg := fmt.Sprintf("%v - %v %v (%v)", c.StatusCode(), c.Method(), c.Path(), time.Since(start))
	log.Info(msg)

	return err
}

type TokenUser struct {
	UserId   string `json:"user_id"`
	Identity string `json:"identity"` // email or username
	Role     string `json:"role"`
}

func EncodeUserMiddleware(c HContext) error {
	log := c.Locals("log").(logger.Logger)

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
	log := c.Locals("log").(logger.Logger)

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
