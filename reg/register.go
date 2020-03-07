package reg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
)

var (
	// Actions ...
	Actions = make(map[string]func(*gin.Context), 0)
	// AuthFunc auth
	AuthFunc = make(map[string]bool)
	// PathActions group by path
	PathActions = make(map[string]string)
	// Ignore ignore auth
	Ignore = make(map[string]func(*gin.Context), 0)
)

// Register register the controllers begin || ignore:需要忽略登录授权的方法
func Register(path string, clt interface{}, auth bool, ignore ...string) {
	if path != "" {
		// Ignore a method of a controller,set unAuth
		var ignoreMap = map[string]string{}
		if len(ignore) > 0 {
			for _, v := range ignore {
				ignoreMap[v] = path
			}
		}
		// filtration
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
		}
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		if _, ok := Actions[path]; ok {
			panic(fmt.Errorf(fmt.Sprintf("The key %s is registered", path)))
			return
		}

		var tr = reflect.TypeOf(clt)
		var num = tr.NumMethod()
		for i := 0; i < num; i++ {
			var name = strings.ToLower(path + tr.Method(i).Name)
			if md, ok := reflect.ValueOf(clt).Method(i).Interface().(func(*gin.Context)); ok {
				if _, iif := ignoreMap[strings.ToLower(tr.Method(i).Name)]; iif && auth { // 存在忽略授权的方法
					Ignore[name] = md
					Actions[name] = md
					continue
				}
				var str = ""
				if strings.Index(name, "default") >= 0 {
					str = path
				} else {
					str = name
				}
				Actions[str] = md
				PathActions[str] = path
				AuthFunc[path] = auth
			} else {
				fmt.Println(fmt.Sprintf("[Registered Routing Info] %s:Method error occurred during conversion,Parameters are not.", name))
			}
		}
	}
}
