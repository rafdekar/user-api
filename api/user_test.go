package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/rafdekar/user-api/db/mock"
	db "github.com/rafdekar/user-api/db/sqlc"
	"github.com/rafdekar/user-api/util"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func randomUser() db.User {
	return db.User{
		ID:        uuid.New(),
		FirstName: util.RandomWord(10),
		LastName:  util.RandomWord(10),
		Email:     util.RandomEmail(),
		Password:  util.RandomWord(10),
		Nickname:  util.RandomWord(10),
		Country:   util.RandomCountry(),
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user *db.User) {
	buffer, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, buffer)

	userToCompare := &db.User{}
	err = json.Unmarshal(buffer, userToCompare)
	require.NoError(t, err)
	require.NotEmpty(t, userToCompare)

	require.Equal(t, user.FirstName, userToCompare.FirstName)
	require.Equal(t, user.LastName, userToCompare.LastName)
	require.Equal(t, user.Email, userToCompare.Email)
	require.Equal(t, user.Nickname, userToCompare.Nickname)
	require.Equal(t, user.Password, userToCompare.Password)
	require.Equal(t, user.Country, userToCompare.Country)
}

func TestCreateUserApi(t *testing.T) {
	user := randomUser()
	dbParams := db.CreateUserParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Nickname:  user.Nickname,
		Country:   user.Country,
	}

	testCases := []struct {
		name          string
		sendEmptyBody bool
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().CreateUser(gomock.Any(), gomock.Eq(dbParams)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, &user)
			},
		},
		{
			name:          "Bad Request",
			sendEmptyBody: true,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().CreateUser(gomock.Any(), gomock.Eq(dbParams)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mockdb.NewMockQuerier(ctrl)
			server := NewServer(querier)

			v.buildStubs(querier)

			var body []byte
			if !v.sendEmptyBody {
				var err error
				body, err = json.Marshal(createUserRequest{
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Email:     user.Email,
					Password:  user.Password,
					Nickname:  user.Nickname,
					Country:   user.Country,
				})
				require.NoError(t, err)
				require.NotEmpty(t, body)
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
			require.NoError(t, err)
			require.NotEmpty(t, req)

			server.router.ServeHTTP(recorder, req)

			v.checkResponse(t, recorder)
		})
	}
}

func TestUpdateUserApi(t *testing.T) {
	user := randomUser()
	dbParams := db.UpdateUserParams{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Nickname:  user.Nickname,
		Country:   user.Country,
	}

	testCases := []struct {
		name          string
		sendEmptyBody bool
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(dbParams)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, &user)
			},
		},
		{
			name:          "Bad Request",
			sendEmptyBody: true,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(dbParams)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(dbParams)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mockdb.NewMockQuerier(ctrl)
			server := NewServer(querier)

			v.buildStubs(querier)

			var body []byte
			if !v.sendEmptyBody {
				var err error
				body, err = json.Marshal(updateUserRequest{
					ID:        user.ID,
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Email:     user.Email,
					Password:  user.Password,
					Nickname:  user.Nickname,
					Country:   user.Country,
				})
				require.NoError(t, err)
				require.NotEmpty(t, body)
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequest("PUT", "/users", bytes.NewBuffer(body))
			require.NoError(t, err)
			require.NotEmpty(t, req)

			server.router.ServeHTTP(recorder, req)

			v.checkResponse(t, recorder)
		})
	}
}

func TestListUsersApi(t *testing.T) {
	n := 2
	users := make([]db.User, n)

	for i := 0; i < n; i++ {
		users[i] = randomUser()
	}

	dbParams := db.ListUsersParams{
		Offset: 2,
		Limit:  2,
	}

	testCases := []struct {
		name          string
		sendEmptyBody bool
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().ListUsers(gomock.Any(), gomock.Eq(dbParams)).
					Times(1).
					Return(users, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

			},
		},
		{
			name:          "Bad Request",
			sendEmptyBody: true,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().ListUsers(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().ListUsers(gomock.Any(), gomock.Eq(dbParams)).
					Times(1).
					Return([]db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mockdb.NewMockQuerier(ctrl)
			server := NewServer(querier)

			v.buildStubs(querier)

			var body []byte
			if !v.sendEmptyBody {
				var err error
				body, err = json.Marshal(listUsersRequest{
					PageSize:   2,
					PageNumber: 2,
				})
				require.NoError(t, err)
				require.NotEmpty(t, body)
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequest("GET", "/users", bytes.NewBuffer(body))
			require.NoError(t, err)
			require.NotEmpty(t, req)

			server.router.ServeHTTP(recorder, req)

			v.checkResponse(t, recorder)
		})
	}
}

func TestDeleteUserApi(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name          string
		sendEmptyBody bool
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "Bad Request",
			sendEmptyBody: true,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mockdb.NewMockQuerier(ctrl)
			server := NewServer(querier)

			v.buildStubs(querier)

			var body []byte
			if !v.sendEmptyBody {
				var err error
				body, err = json.Marshal(deleteUserRequest{
					ID: user.ID,
				})
				require.NoError(t, err)
				require.NotEmpty(t, body)
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequest("DELETE", "/users", bytes.NewBuffer(body))
			require.NoError(t, err)
			require.NotEmpty(t, req)

			server.router.ServeHTTP(recorder, req)

			v.checkResponse(t, recorder)
		})
	}
}
