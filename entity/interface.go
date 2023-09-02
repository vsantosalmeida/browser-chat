package entity

type AuthenticatedUser interface {
	GetId() int
	GetUsername() string
}
