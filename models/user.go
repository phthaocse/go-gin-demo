package models

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/phthaocse/go-gin-demo/schema"
	"github.com/phthaocse/go-gin-demo/utils"
	"log"
	"time"
)

type User struct {
	Id        int        `db:"pk" json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (u *User) CreateFrom(row pgx.Row) error {
	err := row.Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.IsActive, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func processMultiRow(rows *sql.Rows) ([]*User, error) {
	res := make([]*User, 0)
	for rows.Next() {
		usr := &User{}
		err := usr.CreateFrom(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, usr)
	}
	return res, nil
}

func (u *User) GetMulti(db *sql.DB, limit, offset int) ([]*User, error) {
	rows, err := db.Query(`SELECT * FROM "user" LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return processMultiRow(rows)
}

func (u *User) GetAll(db *sql.DB) ([]*User, error) {
	rows, err := db.Query(`SELECT * FROM "user"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return processMultiRow(rows)
}

func (u *User) GetByPk(db *sql.DB) error {
	row, err := GetByPK(u, db)
	if err != nil {
		return err
	}
	err = u.CreateFrom(row)
	if err != nil {
		return err
	}
	u.Password = ""
	return nil
}

func (u *User) GetByEmail(db *sql.DB) (*User, error) {
	row := db.QueryRow(`SELECT * FROM "user" WHERE email = $1`, u.Email)
	if row.Err() != nil {
		return nil, row.Err()
	}
	err := u.CreateFrom(row)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) GetByUsername(db *sql.DB) (*User, error) {
	row := db.QueryRow(`SELECT * FROM "user" WHERE username = $1`, u.Username)
	if row.Err() != nil {
		return nil, row.Err()
	}
	err := u.CreateFrom(row)
	if err != nil {
		return nil, err
	}
	return u, nil
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
	err = db.QueryRow(query, u.Username, u.Email, hashPwd, schema.MemberRole.String()).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (u *User) UpdateActiveStatus(db *sql.DB, status bool) error {
	query := `UPDATE "user" 
    	SET is_active = $1
    	WHERE id = $2
		RETURNING *`
	row := db.QueryRow(query, status, u.Id)
	err := u.CreateFrom(row)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) IsExist(db *sql.DB) bool {
	user, err := u.GetByEmail(db)
	if err != nil {
		log.Println(err)
		return false
	}
	return user != nil
}

func (u *User) IsUsernameExisted(db *sql.DB) bool {
	user, err := u.GetByUsername(db)
	if err != nil {
		log.Println(err)
		return false
	}
	return user != nil
}
