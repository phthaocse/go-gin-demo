package test

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	"github.com/phthaocse/go-gin-demo/server"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var svr *server.Server

func clearDB() {

}
func TestMain(m *testing.M) {
	config := server.GetSrvConfig()
	svr = server.CreateServer(config)
	code := m.Run()
	clearDB()
	os.Exit(code)
}

func TestRegisterBadRequest(t *testing.T) {
	json := []byte(`{"email": "thao@email.com"}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(json))
	svr.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
