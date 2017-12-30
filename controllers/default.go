package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	url := "/admin/index/index"
	c.Redirect(url, 302)
}
