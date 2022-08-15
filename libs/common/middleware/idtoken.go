package middleware

import (
	"encoding/base64"
	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/common/logger"
)

type TokenUser struct {
	UserId   string `json:"user_id"`
	Identity string `json:"identity"` // email or username
	Role     string `json:"role"`
}

func EncodeUserMiddleware(c common.HContext) error {
	log := c.Locals("log").(logger.Interface)

	idToken := c.Get("X-Id-Token")

	if idToken != "" {
		return common.ResponseError(c, common.ErrNotAllowIdToken)
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
		return common.ResponseError(c, common.ErrInvalidIdToken)
	}

	idToken = base64Encode(string(jsonStr))

	c.Set("X-Id-Token", idToken)

	return c.Next()
}

func DecodeUserMiddleware(c common.HContext) error {
	public := c.Locals("public").(bool)
	if public {
		return c.Next()
	}
	idToken := c.Get("X-Id-Token")

	if idToken == "" {
		return common.ResponseError(c, common.ErrNoIdToken)
	}

	jsonStr, ok := base64Decode(idToken)
	if !ok {
		return common.ResponseError(c, common.ErrInvalidIdToken)
	}

	tu := TokenUser{}
	err := json.Unmarshal([]byte(jsonStr), &tu)
	if err != nil {
		log := c.Locals("log").(logger.Interface)
		log.Error(err.Error())
		return common.ResponseError(c, common.ErrInvalidIdToken)
	}

	c.Locals("user", tu)

	return c.Next()
}

func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64Decode(str string) (string, bool) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", false
	}
	return string(data), true
}
