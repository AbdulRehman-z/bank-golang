package api

import (
	"database/sql"

	"github.com/AbdulRehman-z/bank-golang/types"
	"github.com/gofiber/fiber/v2"
)

func (server *Server) renewAccessTokenHandler(c *fiber.Ctx) error {
	var req types.RenewAccessTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	refreshTokenPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	session, err := server.store.GetSession(c.Context(), refreshTokenPayload.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, INTERNAL_SERVER_ERROR)
	}

	if session.IsBlocked {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	if session.Username != refreshTokenPayload.Username {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshTokenPayload.Username, server.config.ACCESS_TOKEN_DURATION)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, INTERNAL_SERVER_ERROR)
	}

	response := types.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": response,
	})
}
