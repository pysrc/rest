package rest

import (
	"net/http"
	"strings"
)

type RoleMapping struct {
	Roles     []string // url模板
	RolesBool []bool   // 如果数组位置对应的是url路由则为true，参数则为false
	Dofunc    func(http.ResponseWriter, *http.Request, map[string]string)
}

func GetUrlParts(url string) []string {
	url = strings.TrimSuffix(url, "/")
	url = strings.TrimPrefix(url, "/")
	return strings.Split(url, "/")
}

func GetUrlRoleParts(role string) ([]string, []bool) {
	var p = GetUrlParts(role)
	var q = make([]bool, len(p))
	if len(p) == 1 && p[0] == "" {
		return p, q
	}
	for i := 0; i < len(p); i++ {
		if p[i][:1] == ":" {
			q[i] = false
			p[i] = p[i][1:]
		} else {
			q[i] = true
		}
	}
	return p, q
}

func Matchs(urlParts []string, roles []string, rolesBool []bool) map[string]string {
	// 判断参数长度
	if len(urlParts) != len(roles) {
		return nil
	}
	// 判断结构是否匹配
	for i, v := range rolesBool {
		if v && urlParts[i] != roles[i] {
			return nil
		}
	}
	// 满足则新建map接收参数
	var res = make(map[string]string)
	for i, v := range rolesBool {
		if !v {
			res[roles[i]] = urlParts[i]
		}
	}
	return res
}

func Match(url string, roles []string, rolesBool []bool) map[string]string {
	return Matchs(GetUrlParts(url), roles, rolesBool)
}

type Router struct {
	mappings  map[string][]*RoleMapping
	validates []func(http.ResponseWriter, *http.Request) bool // 操作验证,必须全通过才进行下一步操作
}

func (self *Router) AddValidate(vali func(http.ResponseWriter, *http.Request) bool) {
	self.validates = append(self.validates, vali)
}

func (self *Router) ErrorRoute(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (self *Router) Route(method, role string, dofunc func(http.ResponseWriter, *http.Request, map[string]string)) {
	if self.mappings == nil {
		self.mappings = make(map[string][]*RoleMapping, 4)
	}
	p, q := GetUrlRoleParts(role)
	self.mappings[method] = append(self.mappings[method], &RoleMapping{p, q, dofunc})
}

func (self *Router) Handle(w http.ResponseWriter, r *http.Request) {
	for _, v := range self.validates { // 验证
		if !v(w, r) {
			return
		}
	}
	rms := self.mappings[r.Method]
	if rms == nil {
		self.ErrorRoute(w, r)
		return
	}
	for _, v := range rms {
		params := Match(r.URL.Path, v.Roles, v.RolesBool)
		if params != nil {
			v.Dofunc(w, r, params)
			return
		}
	}
	self.ErrorRoute(w, r)
}

func (self *Router) Run(addr string) {
	http.HandleFunc("/", self.Handle)
	http.ListenAndServe(addr, nil)
}
