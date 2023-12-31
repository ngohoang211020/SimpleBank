package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ngohoang211020/simplebank/util"
	"github.com/stretchr/testify/require"
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
	require.True(t, user.PasswordChangedAt.UTC().IsZero())
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

	require.WithinDuration(t, account2.CreatedAt.UTC(), account1.CreatedAt.UTC(), time.Millisecond)
	require.WithinDuration(t, account2.PasswordChangedAt.UTC(), account1.PasswordChangedAt.UTC(), time.Millisecond)

}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
}
