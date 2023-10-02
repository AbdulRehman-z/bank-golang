package api

import (
	"strings"

	"github.com/AbdulRehman-z/bank-golang/token"
	"github.com/gofiber/fiber/v2"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := c.GetReqHeaders()["Authorization"]
		if len(authorizationHeader) == 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		authorizationHeaderSplit := strings.Split(authorizationHeader, " ")
		if len(authorizationHeaderSplit) != 2 {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		authorizationType := strings.ToLower(authorizationHeaderSplit[0])
		if authorizationType != authorizationTypeBearer {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		accessToken := authorizationHeaderSplit[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		c.Locals(authorizationPayloadKey, payload)
		return c.Next()
	}
}
