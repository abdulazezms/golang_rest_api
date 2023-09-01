package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"tutorial.sqlc.dev/app/util"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)

	require.Equal(t, payload.Username, username)
	require.Equal(t, payload.Issuer, "localhost:8080")
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiresAt, expiresAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Second

	token, err := maker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.ErrorContains(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	maker1, err1 := NewPasetoMaker(util.RandomString(32))
	maker2, err2 := NewPasetoMaker(util.RandomString(32))

	require.NoError(t, err1)
	require.NoError(t, err2)

	username := util.RandomOwner()
	duration := time.Minute

	token1, err1 := maker1.CreateToken(username, duration)
	token2, err2 := maker2.CreateToken(username, duration)

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NotEmpty(t, token1)
	require.NotEmpty(t, token2)


	payload1, err1 := maker1.VerifyToken(token2)
	require.EqualError(t, err1, ErrInvalidToken.Error())
	require.Empty(t, payload1)
	require.Zero(t, payload1)

	payload2, err2 := maker2.VerifyToken(token1)
	require.EqualError(t, err2, ErrInvalidToken.Error())
	require.Empty(t, payload2)
	require.Zero(t, payload2)

}
