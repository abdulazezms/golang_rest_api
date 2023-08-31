package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"tutorial.sqlc.dev/app/util"
)

func createTestUser(arg CreateUserParams) (User, error) {
	user, err := testQueries.CreateUser(context.Background(), arg)
	return user, err
}

func getRandomUserParams() CreateUserParams {
	return CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomFullName(),
		Email:          util.RandomEmail(),
	}
}
func TestCreateUser(t *testing.T) {
	arg := getRandomUserParams()
	user, err := createTestUser(arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

}

func TestGetUser(t *testing.T) {
	arg := getRandomUserParams()
	createTestUser(arg)
	retrieved, err := testQueries.GetUser(context.Background(), arg.Username)
	require.NoError(t, err)
	require.NotEmpty(t, retrieved)
	require.Equal(t, arg.Username, retrieved.Username)
	require.Equal(t, arg.HashedPassword, retrieved.HashedPassword)
	require.Equal(t, arg.FullName, retrieved.FullName)
	require.Equal(t, arg.Email, retrieved.Email)
	require.NotZero(t, retrieved.CreatedAt)
	require.True(t, retrieved.PasswordChangedAt.IsZero())
}
