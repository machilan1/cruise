package authdb

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/auth"
)

type dbUser struct {
	ID          int       `db:"user_id"`
	UserType    string    `db:"user_type"`
	DisplayName string    `db:"display_name"`
	Username    string    `db:"username"`
	Email       *string   `db:"email"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func toDBUser(usr auth.User) dbUser {
	return dbUser{
		ID:          usr.ID,
		UserType:    string(usr.UserType),
		DisplayName: usr.DisplayName,
		Username:    usr.Username.String(),
		Email:       usr.Email,
		CreatedAt:   usr.CreatedAt,
		UpdatedAt:   usr.UpdatedAt,
	}
}

type dbPassword struct {
	UserID       int    `db:"user_id"`
	PasswordHash []byte `db:"password"`
}

func toDBPassword(userID int, passwordHash []byte) dbPassword {
	return dbPassword{
		UserID:       userID,
		PasswordHash: passwordHash,
	}
}

func toCoreUser(dbUsr dbUser) (auth.User, error) {
	usrname, err := auth.ParseUsername(dbUsr.Username)
	if err != nil {
		return auth.User{}, err
	}
	userType, err := auth.ParseUserType(dbUsr.UserType)
	if err != nil {
		return auth.User{}, err
	}
	return auth.User{
		ID:          dbUsr.ID,
		UserType:    userType,
		DisplayName: dbUsr.DisplayName,
		Username:    usrname,
		Email:       dbUsr.Email,
		CreatedAt:   dbUsr.CreatedAt,
		UpdatedAt:   dbUsr.UpdatedAt,
	}, nil
}

type dbUserWithPassword struct {
	dbUser
	PasswordHash []byte `db:"password"`
}

func toCoreUserWithPassword(dbUsrWithPW dbUserWithPassword) (auth.UserWithPassword, error) {
	u, err := toCoreUser(dbUsrWithPW.dbUser)
	if err != nil {
		return auth.UserWithPassword{}, err
	}

	return auth.UserWithPassword{
		User:         u,
		PasswordHash: dbUsrWithPW.PasswordHash,
	}, nil
}
