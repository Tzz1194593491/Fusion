package fweb

import (
	"net/http"
	"strings"
)

// 用于存放路由相关方法

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// addRouter 添加路由
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	// 同handle中的key
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// getRoute 获取路由
// 例如：trie树中有一个叶子节点leaf = node{ pattern : /a/:b/c, part : c, children : nil, isWild : true }，
// GET,/a/123/c -> leaf,params { b : 123 }
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	// 用于存放解析出来的路径参数
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	node := root.search(searchParts, 0)
	if node != nil {
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
			}
		}
		return node, params
	}
	return nil, nil
}

// handle 处理请求
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		// 同addRouter中的key
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
