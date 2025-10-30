package util

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

type ctxKey string

func ContextWithFiberLocals(c fiber.Ctx) context.Context {
	ctx := c.Context()
	c.RequestCtx().VisitUserValues(func(k []byte, v any) {
		ctx = context.WithValue(ctx, ctxKey(k), v)
	})
	return ctx
}
