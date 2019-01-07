package rest

import (
	"net/http"
	"strings"
)

type RoleMapping struct {
	Roles     []string
	RolesBool []bool
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
	for i := 0; i < len(p); i++ {
		if p[i][:1] == ":" {
			q[i] = true
			p[i] = p[i][1:]
		} else {
			q[i] = false
		}
	}
	return p, q
}

func Matchs(urlParts []string, roles []string, rolesBool []bool) map[string]string {
	if len(urlParts) != len(roles) {
		return nil
	}
	var res = make(map[string]string)
	for i, v := range rolesBool {
		if v {
			res[roles[i]] = urlParts[i]
		} else {
			if urlParts[i] != roles[i] {
				return nil
			}
		}
	}
	return res
}

func Match(url string, roles []string, rolesBool []bool) map[string]string {
	return Matchs(GetUrlParts(url), roles, rolesBool)
}

type Router struct {
	Mappings map[string][]*RoleMapping
	Validate func(http.ResponseWriter, *http.Request) bool // 操作验证
}

func (self *Router) ErrorRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not found !"))
}

func (self *Router) Route(method, role string, dofunc func(http.ResponseWriter, *http.Request, map[string]string)) {
	if self.Mappings == nil {
		self.Mappings = make(map[string][]*RoleMapping, 4)
	}
	p, q := GetUrlRoleParts(role)
	self.Mappings[method] = append(self.Mappings[method], &RoleMapping{p, q, dofunc})
}

func (self *Router) Handle(w http.ResponseWriter, r *http.Request) {
	if self.Validate != nil {
		if !self.Validate(w, r) { // 验证失败
			self.ErrorRoute(w, r)
			return
		}
	}
	rms := self.Mappings[r.Method]
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
