package middleware

import (
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/common/logger"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/util"
	"golang.org/x/exp/slices"
)

var (
	ErrNoToken      = common.NewUnauthorizedError("the token is required")
	ErrInvalidToken = common.NewUnauthorizedError("the token is invalid or expired")
)

func Authentication(secretKey string, excludeList map[string][]string) common.HandleFunc {
	return func(c common.HContext) error {
		public := false

		if methods, ok := excludeList[c.Path()]; ok {
			public = slices.Contains(methods, c.Method())
		}

		if !public && strings.Contains(c.Path(), "/healthz") {
			public = true
		}

		if !public && strings.Contains(c.Path(), "/swagger/") {
			public = true
		}

		if !public && strings.Contains(c.Path(), "/thirdpartySwagger/") {
			public = true
		}

		if !public {
			auth := c.Authorization()
			// validate token
			if auth == "" {
				return common.ResponseError(c, ErrNoToken)
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			claims, err := util.ValidateToken(token, secretKey)

			if err != nil {
				logger.Default.Error(err.Error())
				return common.ResponseError(c, ErrInvalidToken)
			}

			c.Locals("user", claims)
		}

		c.Locals("public", public)

		return c.Next()
	}
}

func AuthenticationCasbin(secretKey string, enforcer *casbin.Enforcer) common.HandleFunc {
	return func(c common.HContext) error {
		public := false

		enforceContext := casbin.NewEnforceContext("2")

		public, err := enforcer.Enforce(enforceContext, c.Path(), c.Method())
		if err != nil {
			logger.Default.Error(err.Error())
			return common.ResponseError(c, ErrEnforce)
		}

		if !public && strings.Contains(c.Path(), "/healthz") {
			public = true
		}

		if !public && strings.Contains(c.Path(), "/swagger/") {
			public = true
		}

		if !public && strings.Contains(c.Path(), "/thirdpartySwagger/") {
			public = true
		}

		if !public {
			auth := c.Authorization()
			// validate token
			if auth == "" {
				return common.ResponseError(c, ErrNoToken)
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			claims, err := util.ValidateToken(token, secretKey)

			if err != nil {
				logger.Default.Debug(err.Error())
				return common.ResponseError(c, ErrInvalidToken)
			}

			c.Locals("user", claims)
		}

		c.Locals("public", public)

		return c.Next()
	}
}
