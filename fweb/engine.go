package fweb

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc 请求处理方法
type HandlerFunc func(*Context)

type (
	RouterGroup struct {
		prefix     string
		middleware []HandlerFunc
		parent     *RouterGroup
		engine     *Engine
	}

	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup // 存储所有的RouterGroup，包括Engine自己的
	}
)

// 初始化相关

// New 创建一个引擎
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group 创建一个新的分组，并只想唯一的engine
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		engine: engine,
		prefix: group.prefix + prefix,
		parent: group,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 路由相关

func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	// 支持同个分组下的公共前缀
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}

// Run 运行服务
func (group *RouterGroup) Run(addr string) (err error) {
	if !strings.HasPrefix(addr, ":") {
		addr = fmt.Sprintf(":%s", addr)
	}
	return http.ListenAndServe(addr, group)
}

// Use 添加中间件到请求处理中
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middleware = append(group.middleware, middlewares...)
}

// 静态资源处理

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	// 去掉请求前缀，例如前缀为/v1/group，那么/v1/group/assert/main.js会被处理为/assert/main.js
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// 检查对应文件是否存在
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

// 实现Handle接口

// ServeHTTP 每次请求过来，都会走这个方法
func (group *RouterGroup) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	// todo：此处可探索是否可以优化
	for _, group := range group.engine.groups {
		// 此步的意思是，将获取对应组的所有父组中间件
		// 例如：v1组中有个product组，product组中有个get路由，完整路由如下：/v1/product/get
		// 那么获取中间件时，需要获取v1组的和product组的，下方代码正是完成此功能
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middleware...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	group.engine.router.handle(c)
}
