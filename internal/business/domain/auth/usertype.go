package auth

import (
	"fmt"
)

type UserType string

const (
	UserTypeSuperAdmin UserType = "superAdmin"
	UserTypeStaff      UserType = "staff"
	UserTypeUser       UserType = "user"

	UserTypeAdmin UserType = "admin"
)

func ParseUserType(userType string) (UserType, error) {
	switch userType {
	case "superAdmin":
		return UserTypeSuperAdmin, nil
	case "staff":
		return UserTypeStaff, nil
	case "admin":
		return UserTypeAdmin, nil
	case "user":
		return UserTypeUser, nil
	default:
		return "", fmt.Errorf("invalid user type: %q", userType)
	}
}

func ParseUpdatableUserType(userType string) (UserType, error) {
	switch userType {
	case "staff":
		return UserTypeStaff, nil
	case "admin":
		return UserTypeAdmin, nil
	case "user":
		return UserTypeUser, nil
	default:
		return "", fmt.Errorf("invalid update user type: %q", userType)
	}
}
