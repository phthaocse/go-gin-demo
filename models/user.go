package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/phthaocse/go-gin-demo/schema"
	"github.com/phthaocse/go-gin-demo/utils"
	"log"
	"reflect"
	"strings"
	"time"
)

type Model interface {
	GetByPk(db *sql.DB) error
}

func GetByPK(model Model, db *sql.DB) (*sql.Row, error) {
	modelRf := reflect.ValueOf(model).Elem()
	mType := modelRf.Type()
	tableName := strings.ToLower(strings.Split(mType.String(), ".")[1])
	var pk reflect.StructField
	for i := 0; i < mType.NumField(); i++ {
		currField := mType.Field(i)
		if val := currField.Tag.Get("db"); val == "pk" {
			pk = currField
			break
		}
	}
	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE %s = $1`, tableName, pk.Name)
	pkVal := reflect.Indirect(modelRf).FieldByName(pk.Name).Interface()
	row := db.QueryRow(query, pkVal)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return row, nil
}

type User struct {
	Id        int          `db:"pk" json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Password  string       `json:"-"`
	Role      string       `json:"role"`
	IsActive  bool         `json:"is_active"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

func (u *User) GetByPk(db *sql.DB) error {
	row, err := GetByPK(u, db)
	if err != nil {
		return err
	}
	err = row.Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.IsActive, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}
	u.Password = ""
	return nil
}

func (u *User) CreateFrom(row *sql.Row) error {
	err := row.Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.IsActive, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}
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
