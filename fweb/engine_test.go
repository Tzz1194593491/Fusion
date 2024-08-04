package fweb

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	r := New()
	v1 := r.Group("/v1")
	v2 := v1.Group("/v2")
	v3 := v2.Group("/v3")
	if v2.prefix != "/v1/v2" {
		t.Fatal("v2 prefix should be /v1/v2")
	}
	if v3.prefix != "/v1/v2/v3" {
		t.Fatal("v2 prefix should be /v1/v2")
	}
}

func TestUse(t *testing.T) {
	engine := Default()
	middleware := engine.middleware
	if len(middleware) == 0 {
		t.Fatal("Use function is error")
	}
}

func TestGet(t *testing.T) {
	engine := Default()
	engine.GET("/v1/hello", nil)
	engine.GET("/v2/hello", nil)
}

func TestPost(t *testing.T) {
	engine := Default()
	engine.POST("/v1/hello", nil)
	engine.POST("/v2/hello", nil)
}

func TestStatic(t *testing.T) {
	engine := Default()
	engine.Static("/fweb", "../fweb")
}

func TestTotal(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	engine := Default()
	engine.Static("/fweb", "../fweb")
	v1 := engine.Group("/v1")
	{
		v1.GET("/hello", func(c *Context) {
			c.JSON(http.StatusOK, H{"path": c.Path})
		})
		v1.POST("/hello", func(c *Context) {
			c.JSON(http.StatusOK, H{"path": c.Path})
		})
	}
	v2 := engine.Group("/v2")
	{
		v2.GET("/hello", func(c *Context) {
			c.JSON(http.StatusOK, H{"path": c.Path})
		})
		v2.POST("/hello", func(c *Context) {
			c.JSON(http.StatusOK, H{"path": c.Path})
		})
	}
	// 启动
	server := httptest.NewServer(engine)
	defer server.Close()
	getTestFunc := func(t *testing.T, url string) {
		response, _ := client.Get(server.URL + url)
		body, _ := io.ReadAll(response.Body)
		var res map[string]string
		_ = json.Unmarshal(body, &res)
		if res["path"] != url {
			t.Fatal("Get " + url + " fail")
		}
	}
	postTestFunc := func(t *testing.T, url string) {
		response, _ := client.Post(server.URL+url, "", nil)
		body, _ := io.ReadAll(response.Body)
		var res map[string]string
		_ = json.Unmarshal(body, &res)
		if res["path"] != url {
			t.Fatal("Post " + url + " fail")
		}
	}
	assertTestFunc := func(t *testing.T, url string) {
		response, _ := client.Get(server.URL + url)
		body, _ := io.ReadAll(response.Body)
		file, _ := os.ReadFile("./go.mod")
		if !reflect.DeepEqual(file, body) {
			t.Fatal("assertTest fail")
		}
	}
	assert404TestFunc := func(t *testing.T, url string) {
		response, _ := client.Get(server.URL + url)
		if response.StatusCode != http.StatusNotFound {
			t.Fatal("assertTest result error")
		}
	}
	// Get请求
	getTestFunc(t, "/v1/hello")
	getTestFunc(t, "/v2/hello")
	// Post请求
	postTestFunc(t, "/v1/hello")
	postTestFunc(t, "/v2/hello")
	// 静态资源
	assertTestFunc(t, "/fweb/go.mod")
	assert404TestFunc(t, "/fweb/dasdhajsdhksahjkagf")
}
