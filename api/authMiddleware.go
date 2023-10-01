package api

import (
	"fmt"
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
			fmt.Printf("authorization header not provided")
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		authorizationHeaderSplit := strings.Split(authorizationHeader, " ")

		if len(authorizationHeaderSplit) != 2 {
			fmt.Printf("incomplete authorization header")
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		authorizationType := strings.ToLower(authorizationHeaderSplit[0])

		if authorizationType != authorizationTypeBearer {
			fmt.Printf("invalid authorization type")
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		accessToken := authorizationHeaderSplit[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			fmt.Printf("invalid access token: %v", err)
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		c.Locals(authorizationPayloadKey, payload)
		fmt.Printf("payload: %v", payload)
		return c.Next()
	}
}
