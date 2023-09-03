package entity

// AuthenticatedUser interface to encapsulate an authenticated user
// with the required fields.
type AuthenticatedUser interface {
	GetId() int
	GetUsername() string
}
