package schema

import "errors"

// UserRole implementation ideal followed https://threedots.tech/post/safer-enums-in-go/
type UserRole struct {
	value string
}

func (ur *UserRole) String() string {
	return ur.value
}

var (
	Unknown    = UserRole{""}
	AdminRole  = UserRole{"admin"}
	MemberRole = UserRole{"member"}
)

func Role(r string) (UserRole, error) {
	switch r {
	case AdminRole.value:
		return AdminRole, nil
	case MemberRole.value:
		return MemberRole, nil
	}
	return Unknown, errors.New("Unknown role: " + r)
}

type UserRegister struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
