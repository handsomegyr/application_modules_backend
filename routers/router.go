package routers

import (
	"application_modules_backend/controllers"
	"application_modules_backend/controllers/system"
	//"fmt"

	"application_modules_backend/library"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.ErrorController(&controllers.ErrorController{})

	var FilterUser = func(ctx *context.Context) {
		ok := !library.Empty(ctx.Input.Session("admin_id"))
		//fmt.Println(adminid, ctx.Request.RequestURI)
		// 如果不是ajax请求的话
		if !ctx.Input.IsAjax() {
			if !ok && ctx.Request.RequestURI != "/admin/index/login" && ctx.Request.RequestURI != "/admin/index/signin" && ctx.Request.RequestURI != "/admin/index/logout" {
				ctx.Redirect(302, "/admin/index/login")
			}
		} else {
			if !ok {
				panic(500)
			}

		}
	}

	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)

	// 用户登陆登出和dashboard页面
	beego.Router("/", &controllers.MainController{})
	beego.Router("/admin", &controllers.MainController{})
	beego.Router("/admin/", &controllers.MainController{})
	beego.Router("/admin/index", &controllers.MainController{})
	beego.Router("/admin/index/", &controllers.MainController{})

	beego.Router("/admin/index/login", &controllers.IndexController{}, "get:Login")
	beego.Router("/admin/index/signin", &controllers.IndexController{}, "post:Signin")
	beego.Router("/admin/index/logout", &controllers.IndexController{}, "get:Logout")
	beego.Router("/admin/index/index", &controllers.IndexController{}, "get:Index")
	// 角色管理
	beego.Router("/admin/system/role/list", &system.RoleController{}, "get:List")
	beego.Router("/admin/system/role/query", &system.RoleController{}, "post:Query")
	beego.Router("/admin/system/role/add", &system.RoleController{}, "get:Add")
	beego.Router("/admin/system/role/insert", &system.RoleController{}, "post:Insert")
	beego.Router("/admin/system/role/edit", &system.RoleController{}, "get:Edit")
	beego.Router("/admin/system/role/update", &system.RoleController{}, "post:Update")
	beego.Router("/admin/system/role/remove", &system.RoleController{}, "post:Remove")
	beego.Router("/admin/system/role/export", &system.RoleController{}, "get:Export")

	// 测试管理
	beego.Router("/admin/system/test/list", &system.TestController{}, "get:List")
	beego.Router("/admin/system/test/query", &system.TestController{}, "post:Query")
	beego.Router("/admin/system/test/add", &system.TestController{}, "get:Add")
	beego.Router("/admin/system/test/insert", &system.TestController{}, "post:Insert")
	beego.Router("/admin/system/test/edit", &system.TestController{}, "get:Edit")
	beego.Router("/admin/system/test/update", &system.TestController{}, "post:Update")
	beego.Router("/admin/system/test/remove", &system.TestController{}, "post:Remove")
	beego.Router("/admin/system/test/export", &system.TestController{}, "get:Export")
}
