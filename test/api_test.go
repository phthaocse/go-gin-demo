package test

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/phthaocse/go-gin-demo/config"
	"github.com/phthaocse/go-gin-demo/models"
	"github.com/phthaocse/go-gin-demo/server"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var svr *server.Server

func prepareData(dbCon *sql.DB) {
	currDir, _ := os.Getwd()
	sqlScript := filepath.Join(currDir, "test_data.sql")
	queries, err := ioutil.ReadFile(sqlScript)
	if err != nil {
		fmt.Println("Can't read the .sql file")
	}
	if res, err := dbCon.Exec(string(queries)); err != nil {
		fmt.Println("Can't execute the query", err)
	} else {
		fmt.Println(res)
	}
}

func TestMain(m *testing.M) {
	config := config.GetSrvConfig()
	svr = server.CreateServer(config)
	prepareData(svr.Db)
	code := m.Run()
	svr.DbTeardown()
	os.Exit(code)
}

func TestRegisterBadRequestExistedEmail(t *testing.T) {
	json := []byte(`{
		"username": "Thao Phan",
		"email": "thao.phan@email.com",
		"password": "12345678"
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(json))
	svr.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error":"User has been existed"}`, w.Body.String())
}

func TestRegisterBadRequestDuplicatedUsername(t *testing.T) {
	json := []byte(`{
		"username": "admin",
		"email": "test@email.com",
		"password": "12345678"
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(json))
	svr.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error":"Username has been existed"}`, w.Body.String())
}

func TestRegisterSuccessfully(t *testing.T) {
	email := "thao.phan2@email.com"
	json := []byte(fmt.Sprintf(`{
		"username": "Thao Phan",
		"email": "%s",
		"password": "12345678"
	}`, email))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(json))
	svr.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `{"message":"Register new user successfully"}`, w.Body.String())

	userInDb := models.User{Email: email}
	_, err := userInDb.GetByEmail(svr.Db)
	assert.Equal(t, err, nil)
}

func TestLoginFailedWrongEmail(t *testing.T) {
	json := []byte(`{
		"email": "noexisted@email.com",
		"password": "12345678"
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(json))
	svr.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"Email or Password incorrect"}`, w.Body.String())
}

func TestLoginFailedWrongPass(t *testing.T) {
	json := []byte(`{
		"email": "thao.phan@email.com",
		"password": "123456789"
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(json))
	svr.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"Email or Password incorrect"}`, w.Body.String())
}

func TestLoginSuccessfully(t *testing.T) {
	json := []byte(`{
		"email": "thao.phan@email.com",
		"password": "12345678"
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(json))
	svr.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
