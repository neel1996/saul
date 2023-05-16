package constants

type UserStatus int

const (
	NoStatus UserStatus = iota
	NewUser
	ExistingUser
)
