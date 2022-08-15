package middleware

import (
	"strings"

	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/common/logger"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/util"
)

var (
	ErrNoToken      = common.NewUnauthorizedError("the token is required")
	ErrInvalidToken = common.NewUnauthorizedError("the token is invalid or expired")
)

func Authentication(secretKey string) common.HandleFunc {
	return func(c common.HContext) error {
		log := c.Locals("log").(logger.Interface)

		public := c.Locals("public").(bool)

		if !public {
			auth := c.Authorization()
			// validate token
			if auth == "" {
				return common.ResponseError(c, ErrNoToken)
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			claims, err := util.ValidateToken(token, secretKey)

			if err != nil {
				log.Error(err.Error())
				return common.ResponseError(c, ErrInvalidToken)
			}

			c.Locals("claims", claims)
		}

		return c.Next()
	}
}
