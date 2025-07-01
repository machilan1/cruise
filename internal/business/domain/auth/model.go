package auth

import (
	"time"
)

type User struct {
	ID          int
	DisplayName string
	Username    Username
	UserType    UserType
	Email       *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserWithPassword struct {
	User
	PasswordHash []byte
}

type NewUser struct {
	Username    Username
	DisplayName string
	Email       *string
	Password    Password
}

type UpdateUserType struct {
	UserID   int
	UserType UserType
}

type UpdateMe struct {
	DisplayName string
}
