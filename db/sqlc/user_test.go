package db

import (
	"context"
	"database/sql"
	"github.com/rafdekar/user-api/util"
	"github.com/stretchr/testify/require"
	"testing"

	_ "github.com/lib/pq"
)

func createTestUser(t *testing.T) *User {
	params := CreateUserParams{
		FirstName: util.RandomWord(5),
		LastName:  util.RandomWord(5),
		Nickname:  util.RandomWord(5),
		Password:  util.RandomWord(5),
		Email:     util.RandomEmail(),
		Country:   util.RandomCountry(),
	}

	user, err := testQueries.CreateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, params.FirstName, user.FirstName)
	require.Equal(t, params.LastName, user.LastName)
	require.Equal(t, params.Nickname, user.Nickname)
	require.Equal(t, params.Password, user.Password)
	require.Equal(t, params.Email, user.Email)
	require.Equal(t, params.Country, user.Country)

	return &user
}

func TestCreateUser(t *testing.T) {
	createTestUser(t)
}

func TestUpdateUser(t *testing.T) {
	testUser := createTestUser(t)

	params := UpdateUserParams{
		ID:        testUser.ID,
		FirstName: util.RandomWord(5),
		LastName:  util.RandomWord(5),
		Nickname:  util.RandomWord(5),
		Password:  util.RandomWord(5),
		Email:     util.RandomEmail(),
		Country:   util.RandomCountry(),
	}

	result, err := testQueries.UpdateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, testUser.ID, result.ID)
	require.Equal(t, params.FirstName, result.FirstName)
	require.Equal(t, params.LastName, result.LastName)
	require.Equal(t, params.Nickname, result.Nickname)
	require.Equal(t, params.Password, result.Password)
	require.Equal(t, params.Email, result.Email)
	require.Equal(t, params.Country, result.Country)
}

func TestListUsers(t *testing.T) {
	n := 10

	for i := 0; i < n; i++ {
		createTestUser(t)
	}

	params := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	result, err := testQueries.ListUsers(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, result, 5)

	for _, v := range result {
		require.NotEmpty(t, v)
	}
}

func TestDeleteUser(t *testing.T) {
	testUser := createTestUser(t)

	err := testQueries.DeleteUser(context.Background(), testUser.ID)
	require.NoError(t, err)

	result, err := testQueries.GetUser(context.Background(), testUser.ID)
	require.Error(t, err, sql.ErrNoRows)
	require.Empty(t, result)
}
