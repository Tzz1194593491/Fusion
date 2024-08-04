package fweb

import (
	"log"
	"time"
)

// Logger 此中间件是用于记录接口调用耗时的
func Logger() HandlerFunc {
	return func(c *Context) {
		// 开始计时
		t := time.Now()
		// 处理请求
		c.Next()
		// 打印耗时
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
