package controllers

import (
	"application_modules_backend/library"
	//"encoding/json"
	//"fmt"
	"strings"

	"application_modules_backend/models/system"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	ModuleName string
}

func (c *BaseController) Prepare() {
	c.ViewPath = beego.AppConfig.String("viewspath")
	c.Data["viewPath"] = c.ViewPath

	c.Data["baseUrl"] = "/"
	c.Data["resourceUrl"] = "/static/backend/metronic.bootstrap/"
	c.Data["commonResourceUrl"] = "/static/common/"

	controllerName, actionName := c.getControllerAndAction()
	c.Data["controllerName"] = controllerName
	c.Data["actionName"] = actionName
	c.Data["moduleName"] = c.ModuleName
	c.Data["auto_redirect"] = false

	c.pagetitle()

	c.breadcrumb()

	c.sidebarmenu()

	c.userlogin()
}

func (c *BaseController) getControllerAndAction() (string, string) {
	controllerName, actionName := c.GetControllerAndAction()
	controllerName = strings.Replace(controllerName, "Controller", "", 1)
	controllerName = strings.ToLower(controllerName)
	actionName = strings.ToLower(actionName)
	//fmt.Println("controllerName:", controllerName, "actionName:", actionName)
	return controllerName, actionName

}
func (c *BaseController) GetUrl(action string) string {
	controllerName, _ := c.getControllerAndAction()
	//fmt.Println("ModuleName:", c.ModuleName)
	return "/" + strings.Join([]string{c.ModuleName, controllerName, action}, "/")
}

func (c *BaseController) breadcrumb() {
	moduleName := library.Strval(c.Data["moduleName"])
	controllerName := library.Strval(c.Data["controllerName"])
	actionName := library.Strval(c.Data["actionName"])

	url1 := "index/index"
	url2 := controllerName + "/" + actionName

	title := ""
	smalltitle := ""

	if moduleName == "admin/system" {
		title = "系统管理"
		if controllerName == "menu" {
			smalltitle += "菜单"
		} else if controllerName == "role" {
			smalltitle += "角色"
		} else if controllerName == "user" {
			smalltitle += "用户"
		}
	} else if moduleName == "admin/lottery" {
		title = "抽奖管理"
		if controllerName == "activity" {
			smalltitle += "活动管理"
		} else if controllerName == "code" {
			smalltitle += "奖品券码管理"
		} else if controllerName == "exchange" {
			smalltitle += "中奖管理"
		} else if controllerName == "identity" {
			smalltitle += "参与者管理"
		} else if controllerName == "limit" {
			smalltitle += "限制管理"
		} else if controllerName == "prize" {
			smalltitle += "奖品管理"
		} else if controllerName == "rule" {
			smalltitle += "概率管理"
		} else if controllerName == "source" {
			smalltitle += "访问来源管理"
		}
	} else if moduleName == "admin/weixin" {
		title = "微信管理"
		if controllerName == "keyword" {
			smalltitle += "关键词管理"
		} else if controllerName == "reply" {
			smalltitle += "回复管理"
		} else if controllerName == "replytype" {
			smalltitle += "回复类型管理"
		} else if controllerName == "menu" {
			smalltitle += "自定义菜单管理"
		} else if controllerName == "menutype" {
			smalltitle += "菜单类型管理"
		} else if controllerName == "application" {
			smalltitle += "应用管理"
		} else if controllerName == "user" {
			smalltitle += "用户管理"
		} else if controllerName == "source" {
			smalltitle += "原始数据管理"
		} else if controllerName == "msgtype" {
			smalltitle += "消息类型管理"
		} else if controllerName == "qrcode" {
			smalltitle += "二维码推广场景管理"
		} else if controllerName == "scene" {
			smalltitle += "二维码场景管理"
		} else if controllerName == "notkeyword" {
			smalltitle += "非关键字管理"
		} else if controllerName == "page" {
			smalltitle += "自定义页面管理"
		} else if controllerName == "gender" {
			smalltitle += "性别管理"
		} else if controllerName == "scripttracking" {
			smalltitle += "执行时间跟踪统计管理"
		} else if controllerName == "callbackurls" {
			smalltitle += "回调地址安全域名管理"
		} else if controllerName == "subscribeuser" {
			smalltitle += "关注用户管理"
		}
	} else {
		if controllerName == "index" && actionName == "index" {
			title = "Dashbord"
			smalltitle = ""
		}
	}

	if actionName == "list" {
		smalltitle += "列表"
	} else if actionName == "add" {
		smalltitle += "追加"
	} else if actionName == "edit" {
		smalltitle += "编辑"
	}
	map1 := make(map[string]string)
	map1["title"] = title
	map1["smalltitle"] = smalltitle
	map1["url1"] = url1
	map1["url2"] = url2
	c.Data["breadcrumb"] = map1
}

func (c *BaseController) pagetitle() {
	map1 := make(map[string]string)
	map1["title"] = "管理中心"
	map1["smalltitle"] = ""
	controllerName := library.Strval(c.Data["controllerName"])
	if controllerName == "menu" {
		map1["smalltitle"] = "菜单管理"
	}
	if controllerName == "role" {
		map1["smalltitle"] = "角色管理"
	}
	if controllerName == "index" {
		map1["smalltitle"] = ""
	}
	c.Data["pagetitle"] = map1
}

func (c *BaseController) sidebarmenu() {

	moduleName := library.Strval(c.Data["moduleName"])
	controllerName := library.Strval(c.Data["controllerName"])

	requestUrl := moduleName + "/" + controllerName
	roleInfo, _ := c.GetSession("roleInfo").(map[string]interface{})

	menu_list, _ := roleInfo["menu_list"].([]string)
	modelMenu := system.NewMenu()
	menus := modelMenu.GetPrivilege(menu_list, requestUrl)

	map1 := make(map[string]interface{})
	map1["requestUrl"] = requestUrl
	map1["menu_list"] = menu_list
	// 排序输出
	map1["menus"] = menus
	map1["sort_keys"] = modelMenu.Sort_keys
	c.Data["sidebarmenu"] = map1
}

func (c *BaseController) userlogin() {
	map1 := make(map[string]interface{})

	map1["admin_id"] = c.GetSession("admin_id")
	map1["admin_name"] = c.GetSession("admin_name")
	c.Data["session"] = map1
}
