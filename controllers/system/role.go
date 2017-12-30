package system

import (
	"application_modules_backend/controllers"
	"application_modules_backend/library"
	"application_modules_backend/models"
	"application_modules_backend/models/system"

	//"github.com/astaxie/beego/validation"
)

type RoleController struct {
	controllers.FormController
}

func (c *RoleController) Prepare() {

	c.ListTplName = "system/role/list.phtml"
	c.AddTplName = "system/role/add.phtml"
	c.EditTplName = "system/role/edit.phtml"

	c.FormName = "角色"
	c.Model1 = system.NewRole()
	c.GetSchemas = c.getSchemas

	c.GetPartials4List = c.getPartials4List
	c.GetList4Show = c.getList4Show
	c.GetData4Add = c.getData4AddOrEdit
	c.GetData4Edit = c.getData4AddOrEdit

	c.Validate4Insert = c.validate4Insert
	c.Validate4Update = c.validate4Update
	c.Validate4Delete = c.validate4Delete

	c.FormController.Prepare()

}

func (c *RoleController) getData4AddOrEdit(row map[string]interface{}) {
	c.Data["privList"] = c.getPrivList(row)
}

func (c *RoleController) getModel() *system.Role {
	modelRole := system.NewRole()
	return modelRole
}

func (c *RoleController) getName() string {
	return "角色"
}

func (c *RoleController) getSchemas(schemas map[string]map[string]interface{}) map[string]map[string]interface{} {

	schemas["name"] = c.CreateSchema("name", "角色名",
		map[string]interface{}{"type": "string", "length": "30"},
		map[string]interface{}{"required": true},
		map[string]interface{}{"input_type": "text", "is_show": true},
		map[string]interface{}{"is_show": true},
		map[string]interface{}{"is_show": true, "placeholder": "角色名..."},
		map[string]interface{}{"is_show": true})

	schemas["alias"] = c.CreateSchema("alias", "角色别名",
		map[string]interface{}{"type": "string", "length": "30"},
		map[string]interface{}{"required": true},
		map[string]interface{}{"input_type": "text", "is_show": true},
		map[string]interface{}{"is_show": true},
		map[string]interface{}{"is_show": true, "placeholder": "角色别名..."},
		map[string]interface{}{"is_show": true})

	schemas["desc"] = c.CreateSchema("desc", "角色描述",
		map[string]interface{}{"type": "string", "length": "100"},
		map[string]interface{}{"required": false},
		map[string]interface{}{"input_type": "textarea", "is_show": true},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false})

	schemas["menu_list"] = c.CreateSchema("menu_list", "菜单设置",
		map[string]interface{}{"type": "array", "length": "4096"},
		map[string]interface{}{"required": false},
		map[string]interface{}{"is_show": true, "partial": "system/role/partials/settings.phtml"},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false})
	//map[string]interface{}{"is_show": true, "input_type": "select", "items": c.GetTrueOrFalseDatas(), "defaultValues": []string{"1"}})

	schemas["operation_list"] = c.CreateSchema("operation_list", "操作设置",
		map[string]interface{}{"type": "array", "length": "4096"},
		map[string]interface{}{"required": false},
		map[string]interface{}{"is_show": true, "partial": "partials/empty.phtml"},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false})

	return schemas
}

func (c *RoleController) getPartials4List() []string {
	return []string{}
	//return []string{"../submodules/post/views/partials/postverify"}
}

func (c *RoleController) getList4Show(input models.Input, list map[string]interface{}) map[string]interface{} {
	// foreach ($list['data'] as &$item) {
	// $item['show_name'] = str_repeat('&nbsp;', $item['level'] * 4) . $item['name'];
	// }
	return list
}

func (c *RoleController) validate4Insert(input models.Input, row map[string]interface{}) bool {
	// do other validation
	c.getModel().CheckName(library.Strval(input.GetKeys("id")), library.Strval(input.GetKeys("name")))
	return true
}

func (c *RoleController) validate4Update(input models.Input, row map[string]interface{}) bool {
	// do other validation
	c.getModel().CheckName(library.Strval(input.GetKeys("id")), library.Strval(input.GetKeys("name")))
	return true
}

func (c *RoleController) validate4Delete(input models.Input, row map[string]interface{}) bool {
	// do other validation
	return true
}

func (c *RoleController) getPrivList(row map[string]interface{}) map[string]interface{} {

	menu_list, _ := row["menu_list"].([]string)
	operation_list, _ := row["operation_list"].([]string)

	//fmt.Println("menu_list:", menu_list)
	/* 获取菜单数据 */
	modelMenu := system.NewMenu()
	menu_priv_arr := modelMenu.GetPrivilege(menu_list, "")
	menuMap := make(map[string]interface{}, 0)
	menuMap["values"] = menu_priv_arr
	menuMap["field"] = "menu_list"

	/* 获取操作数据 */
	modelResource := system.NewResource()
	operation_priv_arr := modelResource.GetPrivilege("admin", operation_list)

	resourceMap := make(map[string]interface{}, 0)
	resourceMap["values"] = operation_priv_arr
	resourceMap["field"] = "operation_list"

	privList := make(map[string]interface{}, 0)
	privList["菜单设置"] = menuMap
	privList["操作设置"] = resourceMap

	return privList
}
