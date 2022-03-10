package models

import (
	"database/sql"
	"errors"
	"github.com/phthaocse/go-gin-demo/utils"
	"log"
	"time"
)

const (
	AdminRole  = "admin"
	MemberRole = "member"
)

type User struct {
	Id        int
	Username  string
	Email     string
	Password  string
	Role      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

func (u *User) GetByEmail(db *sql.DB) error {
	row := db.QueryRow(`SELECT * FROM "user" WHERE email = $1`, u.Email)
	if row.Err() != nil {
		return row.Err()
	}
	err := row.Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.IsActive, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Create(db *sql.DB) (int, error) {
	var userId int
	query := `INSERT INTO "user" 
    	(username, email, password, role)
		VALUES($1, $2, $3, $4) RETURNING id`
	hashPwd, err := utils.HashPassword(u.Password)
	if err != nil {
		return 0, errors.New("can't hash password")
	}
	err = db.QueryRow(query, u.Username, u.Email, hashPwd, MemberRole).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (u *User) IsExist(db *sql.DB) bool {
	err := u.GetByEmail(db)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
