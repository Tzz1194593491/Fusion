package fweb

import (
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	context := newContext(w, r)
	logger := Logger()
	logger(context)
}
