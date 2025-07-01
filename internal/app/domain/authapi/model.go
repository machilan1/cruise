package authapi

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/framework/validate"
)

type Me struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	UserType    string    `json:"userType"`
	DisplayName string    `json:"displayName"`
	Email       *string   `json:"email"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func toAppMe(usr auth.User) Me {
	return Me{
		ID:          usr.ID,
		Username:    usr.Username.String(),
		DisplayName: usr.DisplayName,
		UserType:    string(usr.UserType),
		Email:       usr.Email,
		CreatedAt:   usr.CreatedAt,
		UpdatedAt:   usr.UpdatedAt,
	}
}

func toAppUsers(usrs []auth.User) []Me {
	users := make([]Me, len(usrs))
	for i, usr := range usrs {
		users[i] = toAppMe(usr)
	}
	return users
}

type AppRegisterInput struct {
	Username    string  `json:"username" validate:"required"`
	DisplayName string  `json:"displayName" validate:"required"`
	Email       *string `json:"email" validate:"omitempty,email"`
	Password    string  `json:"password" validate:"required"`
}

func (app AppRegisterInput) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

func toCoreNewUser(app AppRegisterInput) (auth.NewUser, error) {
	usrname, err := auth.ParseUsername(app.Username)
	if err != nil {
		return auth.NewUser{}, err
	}

	password, err := auth.ParsePassword(app.Password)
	if err != nil {
		return auth.NewUser{}, err
	}

	return auth.NewUser{
		Username:    usrname,
		DisplayName: app.DisplayName,
		Email:       app.Email,
		Password:    password,
	}, nil
}

type AppLoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (app AppLoginInput) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

type AppResetPasswordInput struct {
	NewPassword string `json:"newPassword" validate:"required,min=8,max=64"`
}

func (app AppResetPasswordInput) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

type AppForgotPasswordInput struct {
	Email string `json:"email" validate:"required"`
}

func (app AppForgotPasswordInput) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

type AppResetPasswordWithTokenInput struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=64"`
}

func (app AppResetPasswordWithTokenInput) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

type AppUpdateUserTypeInput struct {
	UserID   int    `json:"userId" validate:"required"`
	UserType string `json:"userType" validate:"required,oneof=staff user admin"` // 不可以新增或刪除superAdmin, admin目前為有設計但備而不用，前端鎖住假裝admin這個type不存在
}

func (app AppUpdateUserTypeInput) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

func toCoreUpdateUserType(app AppUpdateUserTypeInput) (auth.UpdateUserType, error) {
	userType, err := auth.ParseUpdatableUserType(app.UserType)
	if err != nil {
		return auth.UpdateUserType{}, err
	}
	return auth.UpdateUserType{
		UserID:   app.UserID,
		UserType: userType,
	}, nil
}

type AppUpdateUserTypeInputs struct {
	UpdateUserTypeUsers []AppUpdateUserTypeInput `json:"updateUserTypeUsers" validate:"required"`
}

func toCoreUpdateUserTypeInputs(app AppUpdateUserTypeInputs) ([]auth.UpdateUserType, error) {
	updateUsers := make([]auth.UpdateUserType, len(app.UpdateUserTypeUsers))
	for i, user := range app.UpdateUserTypeUsers {
		updateUser, err := toCoreUpdateUserType(user)
		if err != nil {
			return nil, err
		}
		updateUsers[i] = updateUser
	}
	return updateUsers, nil
}

type AppUpdateMe struct {
	DisplayName string `json:"displayName" validate:"required"`
}

func toCoreUpdateMe(app AppUpdateMe) auth.UpdateMe {
	return auth.UpdateMe{
		DisplayName: app.DisplayName,
	}
}
