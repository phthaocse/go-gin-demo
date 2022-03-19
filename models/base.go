package models

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
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
