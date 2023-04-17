package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocks "github.com/flukis/simplebank/db/mock"
	db "github.com/flukis/simplebank/db/sqlc"
	"github.com/flukis/simplebank/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T, pwd string) db.User {
	id, err := uuid.NewUUID()
	require.NoError(t, err)
	return db.User{
		ID:             id,
		Username:       util.GenRandomOwner(),
		FullName:       util.GenRandomOwner(),
		HashedPassword: pwd,
		Email:          util.GenRandomEmail(),
	}
}

func TestCreateUserAPI(t *testing.T) {
	arg := util.Argon2Param{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	pwd := "wap12345"
	hashedPwd, err := util.GenerateHashFromPassword(pwd, arg)
	require.NoError(t, err)

	user := randomUser(t, hashedPwd)

	testCases := []struct {
		name  string
		body  any
		build func(store *mocks.Store)
		check func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			body: createUserRequest{
				Email:    user.Email,
				Username: user.Username,
				Fullname: user.FullName,
				Password: pwd,
			},
			build: func(store *mocks.Store) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					FullName:       user.FullName,
					HashedPassword: user.HashedPassword,
					Email:          user.Email,
				}
				store.On("CreateUser", mock.Anything, mock.MatchedBy(func(q db.CreateUserParams) bool {
					if q.Email != arg.Email {
						return false
					}
					if q.FullName != arg.FullName {
						return false
					}
					if q.Username != arg.Username {
						return false
					}
					return true
				})).
					Return(user, nil).
					Once()
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchUser(t, rec.Body, createUserSuccessResponse{Data: generateUserResponse(user)})
			},
		},
		{
			name: "StatusBadRequestPasswordMin8Word",
			body: createUserRequest{
				Email:    user.Email,
				Username: user.Username,
				Fullname: user.FullName,
				Password: "bbb",
			},
			build: func(store *mocks.Store) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					FullName:       user.FullName,
					HashedPassword: user.HashedPassword,
					Email:          user.Email,
				}
				store.On("CreateUser", mock.Anything, mock.MatchedBy(func(q db.CreateUserParams) bool {
					if q.Email != arg.Email {
						return false
					}
					if q.FullName != arg.FullName {
						return false
					}
					if q.Username != arg.Username {
						return false
					}
					return true
				})).
					Return(user, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
	}

	for _, ts := range testCases {
		t.Run(ts.name, func(t *testing.T) {
			store := &mocks.Store{}
			ts.build(store)

			dur, err := time.ParseDuration("1m")
			require.NoError(t, err)

			server, err := NewServer(store, util.Config{
				TokenSymetricKey:    "12345678901234567890123456789012",
				AccessTokenDuration: dur,
			})
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			data, err := json.Marshal(ts.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(data))
			require.NoError(t, err)
			req.Header = http.Header{
				"Content-Type": {"application/json"},
			}

			server.router.ServeHTTP(rec, req)
			ts.check(t, rec)
		})
	}
}

func requireBodyMatchUser[V createUserSuccessResponse](t *testing.T, body *bytes.Buffer, res V) {
	bodyData, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotData V
	err = json.Unmarshal(bodyData, &gotData)
	require.NoError(t, err)
	require.Equal(t, res, gotData)
}
