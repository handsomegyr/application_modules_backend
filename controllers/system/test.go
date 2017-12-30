package system

import (
	"application_modules_backend/controllers"
	"application_modules_backend/library"
	"application_modules_backend/models"
	"application_modules_backend/models/system"

	//"github.com/astaxie/beego/validation"
)

type TestController struct {
	controllers.FormController
}

func (c *TestController) Prepare() {

	c.ListTplName = "system/test/list.phtml"
	c.AddTplName = "system/test/add.phtml"
	c.EditTplName = "system/test/edit.phtml"

	c.FormName = c.getName()
	c.Model1 = system.NewTest()
	c.GetSchemas = c.getSchemas

	c.GetPartials4List = c.getPartials4List
	c.GetList4Show = c.getList4Show

	c.Validate4Insert = c.validate4Insert
	c.Validate4Update = c.validate4Update
	c.Validate4Delete = c.validate4Delete

	c.FormController.Prepare()

}

func (c *TestController) getModel() *system.Test {
	modelTest := system.NewTest()
	return modelTest
}

func (c *TestController) getName() string {
	return "测试"
}

func (c *TestController) getSchemas(schemas map[string]map[string]interface{}) map[string]map[string]interface{} {

	schemas["key"] = c.CreateSchema("key", "关键字",
		map[string]interface{}{"type": "string", "length": "30"},
		map[string]interface{}{"required": true},
		map[string]interface{}{"input_type": "text", "is_show": true},
		map[string]interface{}{"is_show": true, "list_data_name": "show_name", "width": "10%", "name": "关键字2"},
		map[string]interface{}{"is_show": true, "placeholder": "关键字..."},
		map[string]interface{}{"is_show": true})

	schemas["json1"] = c.CreateSchema("json1", "JSON配置",
		map[string]interface{}{"type": "json", "length": "1024"},
		map[string]interface{}{"required": true},
		map[string]interface{}{"input_type": "textarea", "is_show": true},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false, "placeholder": "JSON配置..."},
		map[string]interface{}{"is_show": false})

	schemas["is_valid"] = c.CreateSchema("is_valid", "是否有效",
		map[string]interface{}{"type": "boolean", "length": "1", "defaultValue": "1"},
		map[string]interface{}{"required": true},
		map[string]interface{}{"input_type": "radio", "is_show": true, "items": c.GetTrueOrFalseDatas()},
		map[string]interface{}{"is_show": true, "list_type": "1"},
		map[string]interface{}{"defaultValues": []string{"1"}, "is_show": true, "placeholder": "是否有效...", "input_type": "select", "items": func() []map[string]interface{} { return c.GetTrueOrFalseDatas() }},
		map[string]interface{}{"is_show": true})

	schemas["valid_time"] = c.CreateSchema("valid_time", "验证时间",
		map[string]interface{}{"type": "datetime", "length": "19", "defaultValue": library.GetCurrentTime()},
		map[string]interface{}{"required": false},
		map[string]interface{}{"input_type": "datetimepicker", "is_show": true},
		map[string]interface{}{"is_show": true, "width": "10%"},
		map[string]interface{}{"is_show": true, "input_type": "datetimepicker", "condition_type": "period"},
		map[string]interface{}{"is_show": false})

	schemas["bargain_period"] = c.CreateSchema("bargain_period", "砍价的时间段(小时)",
		map[string]interface{}{"type": "integer", "length": "10"},
		map[string]interface{}{"required": true},
		map[string]interface{}{"input_type": "number", "is_show": true},
		map[string]interface{}{"is_show": true},
		map[string]interface{}{"is_show": true, "input_type": "number", "condition_type": "period"},
		map[string]interface{}{"is_show": true})

	return schemas
}

func (c *TestController) getPartials4List() []string {
	return []string{}
	//return []string{"../submodules/post/views/partials/postverify"}
}

func (c *TestController) getList4Show(input models.Input, list map[string]interface{}) map[string]interface{} {
	// foreach ($list['data'] as &$item) {
	// $item['show_name'] = str_repeat('&nbsp;', $item['level'] * 4) . $item['name'];
	// }
	slice1, _ := list["data"].([]map[string]interface{})
	for key, item := range slice1 {
		item["show_name"] = "KEY:" + library.Strval(item["key"])
		slice1[key] = item
	}
	list["data"] = slice1
	return list
}

func (c *TestController) validate4Insert(input models.Input, row map[string]interface{}) bool {
	// do other validation
	c.getModel().CheckKey(library.Strval(input.GetKeys("id")), library.Strval(input.GetKeys("key")))
	return true
}

func (c *TestController) validate4Update(input models.Input, row map[string]interface{}) bool {
	// do other validation
	c.getModel().CheckKey(library.Strval(input.GetKeys("id")), library.Strval(input.GetKeys("key")))
	return true
}

func (c *TestController) validate4Delete(input models.Input, row map[string]interface{}) bool {
	// do other validation
	return true
}
