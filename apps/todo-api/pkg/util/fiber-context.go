package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/somprasongd/go-monorepo/common"
)

type fiberContext struct {
	*fiber.Ctx
}

func newFiberContext(c *fiber.Ctx) common.HContext {
	return &fiberContext{
		Ctx: c,
	}
}

func (c *fiberContext) Get(key string, defaultValue ...string) string {
	return c.Ctx.Get(key, defaultValue...)
}

func (c *fiberContext) Set(key string, value string) {
	c.Ctx.Set(key, value)
}

func (c *fiberContext) Authorization() string {
	return c.Get("Authorization")
}

func (c *fiberContext) RequestId() string {
	return c.Ctx.GetRespHeader("X-Request-ID")
}

func (c *fiberContext) IP() string {
	return c.Ctx.IP()
}

func (c *fiberContext) Port() string {
	return c.Ctx.Port()
}

func (c *fiberContext) Method() string {
	return c.Ctx.Method()
}

func (c *fiberContext) Path() string {
	return c.Ctx.Path()
}

func (c *fiberContext) BodyParser(v interface{}) error {
	return c.Ctx.BodyParser(v)
}

func (c *fiberContext) QueryParser(v interface{}) error {
	return c.Ctx.QueryParser(v)
}

func (c *fiberContext) Query(key string, defaultValue ...string) string {
	return c.Ctx.Query(key, defaultValue...)
}

func (c *fiberContext) Params(key string, defaultValue ...string) string {
	return c.Ctx.Params(key, defaultValue...)
}

func (c *fiberContext) Locals(key string, value ...interface{}) interface{} {
	return c.Ctx.Locals(key, value...)
}

func (c *fiberContext) Next() error {
	return c.Ctx.Next()
}

func (c *fiberContext) SendStatus(code int) error {
	return c.Ctx.SendStatus(code)
}

func (c *fiberContext) SendJSON(code int, data interface{}) error {
	c.Ctx.Status(code)
	return c.Ctx.JSON(data)
}

func (c *fiberContext) StatusCode() int {
	return c.Ctx.Response().StatusCode()
}

func WrapFiberHandler(h common.HandleFunc) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return h(newFiberContext(c))
	}
}
