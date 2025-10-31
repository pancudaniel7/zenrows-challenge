package util

import (
    "context"

    "github.com/gofiber/fiber/v3"
)

// ctxKey defines a private type for context keys to avoid collisions.
type ctxKey string

// CtxAuthUserIDKey is the context key holding the authenticated user ID.
// Use ctx.Value(util.CtxAuthUserIDKey) to retrieve it as a string.
const CtxAuthUserIDKey ctxKey = "auth_user_id"

// ContextWithFiberLocals returns a stdlib context that includes selected
// values from Fiber locals.
func ContextWithFiberLocals(c fiber.Ctx) context.Context {
    ctx := context.Background()

    if v := c.Locals("auth_user_id"); v != nil {
        ctx = context.WithValue(ctx, CtxAuthUserIDKey, v)
    }

    return ctx
}
