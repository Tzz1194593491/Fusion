package fweb

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func createNewContext(method string, path string, w *httptest.ResponseRecorder, r *http.Request) (*Context, *httptest.ResponseRecorder, *http.Request) {
	if w == nil {
		w = httptest.NewRecorder()
	}
	if r == nil {
		r = httptest.NewRequest(method, path, nil)
	}
	return newContext(w, r), w, r
}

func TestNewContext(t *testing.T) {
	createNewContext("GET", "/test", nil, nil)
}

func TestContext_Fail(t *testing.T) {
	code := http.StatusInternalServerError
	msg := "error"
	exp := struct {
		Message string `json:"message"`
	}{Message: msg}
	expJson, _ := json.Marshal(exp)
	context, w, _ := createNewContext("GET", "/test", nil, nil)
	context.Fail(code, msg)
	assert.Equal(t, code, w.Code)
	assert.Equal(t, string(expJson), strings.ReplaceAll(w.Body.String(), "\n", ""))
}

func TestContext_JSON(t *testing.T) {
	code := http.StatusInternalServerError
	msg := "error"
	exp := struct {
		Message string `json:"message"`
	}{Message: msg}
	expJson, _ := json.Marshal(exp)
	context, w, _ := createNewContext("GET", "/test", nil, nil)
	context.JSON(code, exp)
	assert.Equal(t, code, w.Code)
	assert.Equal(t, string(expJson), strings.ReplaceAll(w.Body.String(), "\n", ""))
}

func TestContext_Data(t *testing.T) {
	code := http.StatusOK
	res := make([]byte, 10)
	context, w, _ := createNewContext("GET", "/test", nil, nil)
	context.Data(code, res)
	assert.Equal(t, code, w.Code)
	assert.Equal(t, res, w.Body.Bytes())
}

func TestContext_HTML(t *testing.T) {
	code := http.StatusOK
	res :=
		"<html>" +
			"<body>" +
			"<h1>hi</h1>" +
			"</body>" +
			"</html>"
	context, w, _ := createNewContext("GET", "/test", nil, nil)
	context.HTML(code, res)
	assert.Equal(t, code, w.Code)
	assert.Equal(t, res, w.Body.String())
}

func TestContext_Query(t *testing.T) {
	username := "123"
	password := "123"
	context, _, _ := createNewContext("GET", fmt.Sprintf("/test?username=%s&password=%s", username, password), nil, nil)
	assert.Equal(t, username, context.Query("username"))
	assert.Equal(t, password, context.Query("password"))
}

func TestContext_PostForm(t *testing.T) {
	username := "123"
	password := "123"
	request := httptest.NewRequest("POST", "/test", nil)
	request.PostForm = url.Values{"username": {username}, "password": {password}}
	context, _, _ := createNewContext("GET", "", nil, request)
	assert.Equal(t, username, context.PostForm("username"))
	assert.Equal(t, password, context.PostForm("password"))
}
