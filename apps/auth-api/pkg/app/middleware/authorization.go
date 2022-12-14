package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/common/logger"
)

var (
	ErrUnauthorize = common.NewForbiddenError("you are not allowed to access this resource")
	ErrEnforce     = common.NewUnexpectedError("error occurred while enforce")
)

func Authorize(enforcer *casbin.Enforcer) common.HandleFunc {
	return func(c common.HContext) error {
		log := c.Locals("log").(logger.Interface)
		public := c.Locals("public").(bool)
		if public {
			return c.Next()
		}

		claims := c.Locals("claims").(jwt.MapClaims)
		role := claims["role"].(string)

		ok, err := enforcer.Enforce(role, c.Path(), c.Method())

		if err != nil {
			log.Error(err.Error())
			return common.ResponseError(c, ErrEnforce)
		}

		if !ok {
			return common.ResponseError(c, ErrUnauthorize)
		}

		return c.Next()
	}
}
