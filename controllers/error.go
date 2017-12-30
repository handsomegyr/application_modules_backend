package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Prepare() {
	c.Data["resourceUrl"] = "/static/backend/metronic.bootstrap/"
	c.Data["commonResourceUrl"] = "/static/common/"
	c.Data["controllerName"], c.Data["actionName"] = c.GetControllerAndAction()
}

func (c *ErrorController) Error401() {
	c.TplName = "error/show401.phtml"
}

func (c *ErrorController) Error404() {
	c.Data["content"] = "page not found"
	c.TplName = "error/show404.phtml"
}

func (c *ErrorController) Error500() {
	c.Data["content"] = "server error"
	c.TplName = "error/show500.phtml"
}

func (c *ErrorController) Message() {
	c.Data["content"] = "server error"
	c.TplName = "error/message.phtml"
}
