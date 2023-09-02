package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	mytoken "tutorial.sqlc.dev/app/token"
	"tutorial.sqlc.dev/app/util"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker mytoken.Maker,
	authorizationType string,
	username string,
	duration time.Duration) {

	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	request.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationType, token))
}

func TestAuthMiddleware(t *testing.T) {

	tokenMaker, err := mytoken.NewPasetoMaker(util.RandomString(32))
	//Let's create a token now.
	require.NoError(t, err)
	token, err := tokenMaker.CreateToken(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	t.Log(token)
	//now we must check the middleware
	testCases := []struct {
		name          string
		setUpAuth     func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "No Authorization",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unsupported Authorization",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker) {
				addAuthorization(t, request, tokenMaker, "unsupported", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name: "Invalid Authorization Format",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker) {
				addAuthorization(t, request, tokenMaker, "- - -", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name: "Expired Token",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker mytoken.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			server, err := NewTestServer(t, nil)
			require.NoError(t, err)

			authPath := "/auth"
			server.router.GET(
				authPath,                          //path
				authMiddleWare(server.tokenMaker), //middleware
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{}) //handler after after middleware
				},
			)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)
			tc.setUpAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
