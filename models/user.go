package models

const (
	AdminRole  = "admin"
	MemberRole = "member"
)

type User struct {
	Id       int
	Username string
	Role     string
}
