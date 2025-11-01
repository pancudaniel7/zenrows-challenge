package middleware

import (
	"encoding/base64"
	"strings"

	"zenrows-challenge/internal/core/port"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

const AuthUserIDKey = "auth_user_id"

func BasicAuthCheckMiddleware(svc port.AuthenticationService, v *validator.Validate) fiber.Handler {
	return func(c fiber.Ctx) error {
		h := c.Get("Authorization")
		if h == "" {
			return fiber.ErrUnauthorized
		}

		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Basic") {
			return fiber.ErrUnauthorized
		}

		b, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return fiber.ErrUnauthorized
		}

		creds := string(b)
		i := strings.IndexByte(creds, ':')
		if i <= 0 {
			return fiber.ErrUnauthorized
		}

		user := creds[:i]
		pass := creds[i+1:]
		if user == "" || pass == "" {
			return fiber.ErrUnauthorized
		}

		userID, err := svc.CheckCredentials(user, pass)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		if userID == "" && v.Var(userID, "required,uuid4") != nil {
			return fiber.ErrUnauthorized
		}
		c.Locals(AuthUserIDKey, userID)
		return c.Next()
	}
}
