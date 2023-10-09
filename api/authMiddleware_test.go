package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/AbdulRehman-z/bank-golang/token"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				accessToken, payload, err := tokenMaker.CreateToken("test", time.Minute)
				require.NotEmpty(t, payload)
				require.NoError(t, err)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name: "NoAuthorizationHeader",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InvalidAuthorizationHeader",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				accessToken, payload, err := tokenMaker.CreateToken("test", time.Minute)
				require.NotEmpty(t, payload)
				require.NoError(t, err)
				request.Header.Set("Authorization", fmt.Sprintf(accessToken))
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InvalidAuthorizationType",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				accessToken, payload, err := tokenMaker.CreateToken("test", time.Minute)
				require.NotEmpty(t, payload)
				require.NoError(t, err)
				request.Header.Set("Authorization", fmt.Sprintf("test %s", accessToken))
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				accessToken, payload, err := tokenMaker.CreateToken("test", -time.Minute)
				require.NotEmpty(t, payload)
				require.NoError(t, err)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			server := NewTestServer(t, nil)
			server.router.Get("/auth", AuthMiddleware(server.tokenMaker), func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusOK).JSON(&fiber.Map{})
			})

			request, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)
			response, err := server.router.Test(request)
			require.NoError(t, err)

			tc.checkResponse(t, response)

		})
	}
}
