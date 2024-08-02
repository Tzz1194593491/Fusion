package fweb

import (
	"net/http"
)

// HandlerFunc 请求处理方法
type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

// New 创建一个引擎
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.router.addRouter("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.router.addRouter("POST", pattern, handler)
}

// Run 运行服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP 每次请求过来，都会走这个方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
