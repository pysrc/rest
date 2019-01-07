package main

import (
	"fmt"
	"github.com/pysrc/rest"
	"net/http"
)

func main() {
	var router rest.Router
	router.Validate = func(w http.ResponseWriter, r *http.Request) bool { // 拦截所有请求，并验证
		if r.URL.Path[:4] == "/api" {
			return true
		}
		return false
	}
	router.Route("GET", "/api/:name/:pwd/index", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Write([]byte(fmt.Sprintln(r.Method, params)))
	})
	router.Route("POST", "/api/:name/:pwd/index", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Write([]byte(fmt.Sprintln(r.Method, params)))
	})
	router.Route("PUT", "/api/:name/:pwd/index", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Write([]byte(fmt.Sprintln(r.Method, params)))
	})
	router.Route("DELETE", "/api/:name/:pwd/index", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Write([]byte(fmt.Sprintln(r.Method, params)))
	})
	router.Run("127.0.0.1:8080")
}
