package fweb

import (
	"fmt"
	"reflect"
	"testing"
)

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
