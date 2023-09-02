package presenter

import "github.com/vsantosalmeida/browser-chat/entity"

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
}

type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserOutput struct {
	ID int `json:"id"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func MapEntityToExternalUsers(users []*entity.User) []*User {
	result := make([]*User, 0)

	for _, user := range users {
		result = append(
			result,
			&User{
				ID:       user.ID,
				Username: user.Username,
			},
		)
	}

	return result
}
