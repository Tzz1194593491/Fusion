package fweb

import (
	"net/http/httptest"
	"testing"
)

func TestRecovery(t *testing.T) {
	recovery := Recovery()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	context := newContext(w, r)
	context.handlers = append(context.handlers, func(c *Context) {
		panic("error")
	})
	recovery(context)
}
