package db

import (
	"context"
	"testing"
	"time"

	"github.com/flukis/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createDummyUser(t *testing.T) User {
	args := CreateUserParams{
		Username:       util.GenRandomOwner(),
		FullName:       util.GenRandomOwner(),
		Email:          util.GenRandomEmail(),
		HashedPassword: "secret",
	}

	user, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, args.Username)
	require.Equal(t, user.FullName, args.FullName)
	require.Equal(t, user.Email, args.Email)
	require.Equal(t, user.HashedPassword, args.HashedPassword)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createDummyUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createDummyUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}
