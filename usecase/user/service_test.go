package user_test

import (
	"testing"
	"time"

	"github.com/vsantosalmeida/browser-chat/entity"
	"github.com/vsantosalmeida/browser-chat/usecase/user"
	"github.com/vsantosalmeida/browser-chat/usecase/user/mocks"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/undefinedlabs/go-mpatch"
	"golang.org/x/crypto/bcrypt"
)

var errDB = errors.New("db error")

func TestService_Authenticate(t *testing.T) {
	var (
		username = "test"
		password = "testing"
		hash     = "$2a$10$/LguUiu0z2YSnulRC5NXDe3lbrnBVyCbYfjfP3xRse8DXlseJoI1G"
		expected = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE1Nzc5MjMyMDAsIklEIjowLCJVc2VybmFtZSI6IiJ9.6exc1Ml3XDfEd9swyVI-bxy57GmExUUK0clofSvrH0k"
	)

	repository := mocks.NewRepository(t)
	svc := user.NewService(repository)

	// bypass time.Now function to set a static date for the JWT token
	timePatch, err := mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	})
	assert.NoError(t, err)
	defer timePatch.Unpatch()

	repository.
		On("FindByUsername", "test").
		Return(&entity.User{Password: hash}, nil).
		Once()

	token, err := svc.Authenticate(username, password)
	assert.NoError(t, err)
	assert.Equal(t, expected, token)
}

func TestService_AuthenticateErrors(t *testing.T) {
	var tt = []struct {
		name     string
		username string
		password string
		mockErr  error
		expected string
	}{
		{
			name:     "When password is wrong; should return error",
			username: "test",
			password: "invalid",
			expected: "could not validate user: invalid password",
		},
		{
			name:     "When could not retrieve user entity; should return error",
			username: "test",
			password: "invalid",
			expected: "could not find user: db error",
			mockErr:  errDB,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var hash = "$2a$10$/LguUiu0z2YSnulRC5NXDe3lbrnBVyCbYfjfP3xRse8DXlseJoI1G"

			repository := mocks.NewRepository(t)
			svc := user.NewService(repository)

			repository.
				On("FindByUsername", "test").
				Return(&entity.User{Password: hash}, tc.mockErr).
				Once()

			token, err := svc.Authenticate(tc.username, tc.password)
			assert.EqualError(t, err, tc.expected)
			assert.Empty(t, token)
		})
	}
}

func TestService_CreateUser(t *testing.T) {
	var (
		username = "test"
		password = "testing"
		expected = 1
		hash     = "$2a$10$/LguUiu0z2YSnulRC5NXDe3lbrnBVyCbYfjfP3xRse8DXlseJoI1G"

		userEntity = &entity.User{
			Username: "test",
			Password: "$2a$10$/LguUiu0z2YSnulRC5NXDe3lbrnBVyCbYfjfP3xRse8DXlseJoI1G",
		}
	)

	repository := mocks.NewRepository(t)
	svc := user.NewService(repository)

	// bypass bcrypt.GenerateFromPassword function to set a static password hash
	cryptoPatch, err := mpatch.PatchMethod(bcrypt.GenerateFromPassword, func([]byte, int) ([]byte, error) {
		return []byte(hash), nil
	})
	assert.NoError(t, err)
	defer cryptoPatch.Unpatch()

	repository.
		On("Create", userEntity).
		Return(1, nil).
		Once()

	id, err := svc.CreateUser(username, password)
	assert.NoError(t, err)
	assert.Equal(t, expected, id)
}

func TestService_CreateUserErrors(t *testing.T) {
	var tt = []struct {
		name     string
		userName string
		password string
		mockErr  error
		expected string
	}{
		{
			name:     "When username is empty; should return error",
			userName: "",
			password: "testing",
			expected: "could not create an user object: invalid entity",
		},
		{
			name:     "When password is empty; should return error",
			userName: "test",
			password: "",
			expected: "could not create an user object: invalid entity",
		},
		{
			name:     "When could not create user on DB; should return error",
			userName: "test",
			password: "testing",
			mockErr:  errDB,
			expected: "could not create user on DB: db error",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				hash = "$2a$10$/LguUiu0z2YSnulRC5NXDe3lbrnBVyCbYfjfP3xRse8DXlseJoI1G"

				userEntity = &entity.User{
					Username: "test",
					Password: "$2a$10$/LguUiu0z2YSnulRC5NXDe3lbrnBVyCbYfjfP3xRse8DXlseJoI1G",
				}
			)

			repository := mocks.NewRepository(t)
			svc := user.NewService(repository)

			// bypass bcrypt.GenerateFromPassword function to set a static password hash
			cryptoPatch, err := mpatch.PatchMethod(bcrypt.GenerateFromPassword, func([]byte, int) ([]byte, error) {
				return []byte(hash), nil
			})
			assert.NoError(t, err)
			defer cryptoPatch.Unpatch()

			repository.
				On("Create", userEntity).
				Return(0, tc.mockErr).
				Maybe()

			id, err := svc.CreateUser(tc.userName, tc.password)
			assert.EqualError(t, err, tc.expected)
			assert.Empty(t, id)
		})
	}
}

func TestService_ListUsers(t *testing.T) {
	var (
		usersList = []*entity.User{
			{
				ID:       1,
				Username: "firstUser",
				Password: "password1",
			},
			{
				ID:       2,
				Username: "secondUser",
				Password: "password2",
			},
		}

		expected = []*entity.User{
			{
				ID:       1,
				Username: "firstUser",
				Password: "password1",
			},
			{
				ID:       2,
				Username: "secondUser",
				Password: "password2",
			},
		}
	)

	repository := mocks.NewRepository(t)
	svc := user.NewService(repository)

	repository.
		On("List").
		Return(usersList, nil).
		Once()

	users, err := svc.ListUsers()
	assert.NoError(t, err)
	assert.Equal(t, expected, users)
}

func TestService_ListUsersError(t *testing.T) {
	var (
		expected = "could not retrieve users list: db error"
	)

	repository := mocks.NewRepository(t)
	svc := user.NewService(repository)

	repository.
		On("List").
		Return(nil, errDB).
		Once()

	users, err := svc.ListUsers()
	assert.EqualError(t, err, expected)
	assert.Empty(t, users)
}
