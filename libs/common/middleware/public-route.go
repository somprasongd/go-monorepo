package middleware

import (
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/common/logger"
	"golang.org/x/exp/slices"
)

var (
	ErrEnforce = common.NewUnexpectedError("error occurred while enforce")
)

type PublicRoute struct {
	Path    string
	Methods []string
}

func (pr *PublicRoute) AllowMethod(method string) bool {
	return slices.Contains(pr.Methods, method)
}

func PublicRouteMiddleware(publicList ...PublicRoute) common.HandleFunc {
	return func(c common.HContext) error {
		public := false

		path := c.Path()

		for _, v := range publicList {
			if v.Path == path {
				public = v.AllowMethod(c.Method())
				break
			}
		}

		if !public && strings.Contains(path, "/healthz") {
			public = true
		}

		if !public && strings.Contains(path, "/swagger/") {
			public = true
		}

		if !public && strings.Contains(path, "/thirdpartySwagger/") {
			public = true
		}

		c.Locals("public", public)

		return c.Next()
	}
}

func PublicRouteMiddlewareCasbin(enforcer *casbin.Enforcer, suffix string) common.HandleFunc {
	return func(c common.HContext) error {
		log := c.Locals("log").(logger.Interface)

		public := false

		enforceContext := casbin.NewEnforceContext(suffix)

		public, err := enforcer.Enforce(enforceContext, c.Path(), c.Method())
		if err != nil {
			log.Error(err.Error())
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

		c.Locals("public", public)

		return c.Next()
	}
}
