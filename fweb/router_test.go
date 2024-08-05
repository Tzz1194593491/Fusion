package fweb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// 单元测试

func newTestRouter() *router {
	r := newRouter()
	// 加入/*filepath后无法正常匹配
	r.addRouter("GET", "/*filepath", nil)
	r.addRouter("GET", "/", nil)
	r.addRouter("GET", "/hello/:name", nil)
	r.addRouter("GET", "/assert/*filepath", nil)
	r.addRouter("POST", "/hello/a/b", nil)
	r.addRouter("POST", "/hi/:name", nil)
	return r
}

func newTestRouter2() *router {
	r := newRouter()
	r.addRouter("GET", "/", nil)
	r.addRouter("GET", "/hello/:name", nil)
	r.addRouter("GET", "/assert/*filepath", nil)
	r.addRouter("GET", "/hello/a/b", nil)
	r.addRouter("GET", "/hi/:name", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/abc"), []string{"p", "abc"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/123")
	if n == nil {
		t.Fatal("not matched /hello/:name")
	}
	if ps == nil {
		t.Fatal("params load fail")
	}
	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}
	if ps["name"] != "123" {
		t.Fatal("name should be equal to '123'")
	}
	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
}

func TestGetRouteURL404(t *testing.T) {
	r := newTestRouter2()
	n, ps := r.getRoute("GET", "/hello/okokda/12331231231")
	assert.Nil(t, n, "expect n is nil")
	assert.Nil(t, ps, "expect ps is nil")
}

func TestGetRouteMethod404(t *testing.T) {
	r := newTestRouter2()
	n, ps := r.getRoute("POST", "/hello/123")
	assert.Nil(t, n, "expect n is nil")
	assert.Nil(t, ps, "expect ps is nil")
}

func TestGetRouteHandle404(t *testing.T) {
	r := newTestRouter2()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/hihi/123", nil)
	r.handle(newContext(w, req))
	assert.Equal(t, http.StatusNotFound, w.Code)
}
