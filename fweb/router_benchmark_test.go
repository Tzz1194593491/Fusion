package fweb

import (
	"testing"
)

// 基准测试

func BenchmarkGetRoute(b *testing.B) {
	r := newTestRouter()
	testFn := func() {
		// 成功匹配
		_, _ = r.getRoute("GET", "/hello/123")
		_, _ = r.getRoute("GET", "/assert/123.js")
		_, _ = r.getRoute("GET", "/hello/a/b")
		_, _ = r.getRoute("GET", "/hi/a")
		// 失败匹配
		_, _ = r.getRoute("GET", "/hello/123/123")
		_, _ = r.getRoute("GET", "/hello/123/123/123")
		_, _ = r.getRoute("GET", "/hello/123/123/123/123/123/123/123/123/123/123/123/123/123/123/123/123/123/123")
	}
	for n := 0; n < b.N; n++ {
		testFn()
	}
}
