package controllers

import (
	//"fmt"
	"application_modules_backend/models/system"
	"log"

	"github.com/astaxie/beego/validation"
)

type IndexController struct {
	BaseController
}

func (c *IndexController) Prepare() {
	c.ModuleName = "admin"
	c.BaseController.Prepare()
}

// dashboard页面
func (c *IndexController) Index() {
	c.Layout = "index.phtml"
	c.TplName = "index/index.phtml"
}

//  用户登陆页面
func (c *IndexController) Login() {
	c.Data["form_act"] = c.GetUrl("signin")
	c.TplName = "index/login.phtml"
}

// 用户登陆
func (c *IndexController) Signin() {
	username := c.Input().Get("username")
	password := c.Input().Get("password")

	valid := validation.Validation{}
	valid.Required(username, "You must enter your username.")
	valid.Required(password, "You must enter your password.")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	} else {
		modelUser := system.NewUser()
		userInfo := modelUser.CheckLogin(username, password)
		if len(userInfo) == 0 {
			panic(401)
		}

		// 登陆处理
		modelUser.Login(userInfo)
		//c.StopRun()
		//fmt.Println(userInfo)

		c.SetSession("admin_id", userInfo["_id"])
		c.SetSession("admin_name", userInfo["username"])

		// 获取角色信息
		if role, ok := userInfo["role"].(string); ok {
			modelRole := system.NewRole()
			roleInfo := modelRole.GetInfoById(role)
			c.SetSession("roleInfo", roleInfo)
		} else {
			c.SetSession("roleInfo", make(map[string]interface{}))
		}

		//c.StopRun()
		url := c.GetUrl("index")
		c.Redirect(url, 302)
	}
}

// 用户登出
func (c *IndexController) Logout() {
	c.DestroySession()
	url := c.GetUrl("login")
	c.Redirect(url, 302)
}
