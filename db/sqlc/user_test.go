package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) Users {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Username:       util.RandomOwner(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.Time.IsZero())
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	account1 := createRandomUser(t)
	account2, err := testStore.GetUser(context.Background(), account1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.FullName, account2.FullName)
	require.Equal(t, account1.Email, account2.Email)
	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, account1.HashedPassword, account2.HashedPassword)

	require.WithinDuration(t, account2.CreatedAt.Time, account1.CreatedAt.Time, time.Millisecond)
	require.WithinDuration(t, account2.PasswordChangedAt.Time, account1.PasswordChangedAt.Time, time.Millisecond)

}
