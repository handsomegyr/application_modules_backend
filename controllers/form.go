package controllers

import (
	"application_modules_backend/library"
	"application_modules_backend/models"
	"application_modules_backend/models/system"
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/validation"
)

type FormController struct {
	BaseController
	// 模板页面
	ListTplName string
	AddTplName  string
	EditTplName string

	//表单名称
	FormName string
	// 表操作对象
	Model1 system.IBase
	// 字段说明列表
	GetSchemas func(map[string]map[string]interface{}) map[string]map[string]interface{}

	// 列表增加修改删除页面显示处理
	GetPartials4List func() []string
	GetList4Show     func(input models.Input, list map[string]interface{}) map[string]interface{}
	GetData4Add      func(row map[string]interface{})
	GetData4Edit     func(row map[string]interface{})
	GetData4Delete   func(row map[string]interface{})

	// 增加修改删除数据的检查处理
	Validate4Insert func(input models.Input, row map[string]interface{}) bool
	Validate4Update func(input models.Input, row map[string]interface{}) bool
	Validate4Delete func(input models.Input, row map[string]interface{}) bool

	Schema_keys []string
}

func (c *FormController) Prepare() {
	c.ModuleName = "admin/system"
	c.BaseController.Prepare()

	c.Layout = "index.phtml"
	c.Data["currentTime"] = library.GetCurrentTime()
	c.Data["formName"] = c.getName()

	schemas := c.getSchemas()
	c.Data["schemas"] = schemas
	c.Data["partials4List"] = c.getPartials4List()

	c.Data["schema_keys"] = c.Schema_keys

	c.buildSearch(schemas)
	c.buildList(schemas)

}

//显示列表页面
func (c *FormController) List() {
	c.TplName = c.ListTplName
	c.Data["defaultSort"] = c.getModel().GetDefaultSort()
}

//排序、分页、查询
func (c *FormController) Query() {

	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println("Query err:", err)
			datas := make(map[string]interface{}, 0)
			datas["draw"] = 1
			datas["error"] = library.Strval(err)
			c.Data["json"] = datas
			c.ServeJSON()
		}
	}()

	input := c.getListFilterInput()
	if !input.IsValid() {
		messageInfo := c.GetValidationMessage(input)
		panic(errors.New(messageInfo))
	}

	// 根据检索条件获取列表
	list := c.getList(input)

	// 将列表数据按照画面要求进行显示
	list = c.getList4Show(input, list)

	datas := make(map[string]interface{}, 0)
	datas["draw"] = input.GetKeys("draw")
	datas["recordsTotal"] = list["record_count"]
	datas["recordsFiltered"] = list["record_count"]
	datas["data"] = list["data"]

	c.Data["json"] = datas
	c.ServeJSON()

}

//数据导出
func (c *FormController) Export() {
	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println("Export err:", err)
			datas := make(map[string]interface{}, 0)
			datas["error"] = library.Strval(err)
			c.Data["json"] = datas
			c.ServeJSON()
		}
	}()

	input := c.getListFilterInput()
	if !input.IsValid() {
		messageInfo := c.GetValidationMessage(input)
		panic(errors.New(messageInfo))
	}

	c.EnableRender = false

	//根据检索条件获取列表
	//fmt.Println("GetSort:", input.GetSort())
	list := c.getModel().GetAllList(input)
	datas := c.export(list)
	_, content := library.ArrayToCVS(datas)

	c.Ctx.Output.Header("Content-Type", "text/csv; charset=utf-8")
	c.Ctx.Output.Header("Content-Disposition", `attachment; filename="`+"test"+".csv")
	//c.Ctx.Output.Header("Content-Length", 300)
	//fmt.Println("content:", content)
	c.Ctx.Output.Body(content)
}

//显示添加页面
func (c *FormController) Add() {
	c.TplName = c.AddTplName

	/* 初始化、取得 信息 */
	row := c.getModel().GetEmptyRow(c.getFilterInput())

	c.Data["row"] = row
	c.Data["form_act"] = c.GetUrl("insert")
	if c.GetData4Add != nil {
		c.GetData4Add(row)
	}
}

//追加
func (c *FormController) Insert() {

	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println("Insert err:", err)
			datas := make(map[string]interface{}, 0)
			datas["error"] = library.Strval(err)
			c.Data["json"] = datas
			c.ServeJSON()
		}
	}()

	input := c.getFilterInput()
	if !input.IsValid() {
		messageInfo := c.GetValidationMessage(input)
		panic(errors.New(messageInfo))
	}

	// 在进行插入处理之前进行检查
	c.validate4Insert(input, nil)

	// insert
	c.insert(input, nil)

	/* 添加链接 */
	link := make([]map[string]interface{}, 0)
	map1 := make(map[string]interface{}, 0)
	map1["text"] = "继续添加" + c.getName()
	map1["href"] = c.GetUrl("add")

	map2 := make(map[string]interface{}, 0)
	map2["text"] = "返回" + c.getName() + "列表"
	map2["href"] = c.GetUrl("list")
	link = append(link, map1, map2)

	c.SysMsg(c.getName()+"添加成功!", 0, link, true)

}

//显示编辑页面
func (c *FormController) Edit() {
	c.TplName = c.EditTplName

	input := c.getFilterInput()
	if !input.IsValid("id") {
		messageInfo := c.GetValidationMessage(input)
		panic(errors.New(messageInfo))
	}

	/* 初始化、取得 信息 */
	row := c.getModel().GetInfoById(library.Strval(input.GetKeys("id")))

	c.Data["row"] = row
	c.Data["form_act"] = c.GetUrl("update")

	if c.GetData4Edit != nil {
		c.GetData4Edit(row)
	}

}

//更新
func (c *FormController) Update() {
	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println("Update err:", err)
			datas := make(map[string]interface{}, 0)
			datas["error"] = library.Strval(err)
			c.Data["json"] = datas
			c.ServeJSON()
		}
	}()

	input := c.getFilterInput()
	if !input.IsValid() {
		messageInfo := c.GetValidationMessage(input)
		panic(errors.New(messageInfo))
	}

	// get exist
	row := c.getModel().GetInfoById(library.Strval(input.GetKeys("id")))
	//fmt.Println("id:", input.GetKeys("id"), "row:", row)

	if library.Empty(row) {
		panic(errors.New("更新的数据为空"))
	}
	// 在进行更新处理之前进行检查
	c.validate4Update(input, row)

	// update
	c.update(input, row)

	/* 添加链接 */
	link := make([]map[string]interface{}, 0)
	map1 := make(map[string]interface{}, 0)
	map1["text"] = "返回" + c.getName() + "列表"
	map1["href"] = c.GetUrl("list")
	link = append(link, map1)

	c.SysMsg(c.getName()+"编辑成功!!", 0, link, true)

}

//删除
func (c *FormController) Remove() {
	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println("Remove err:", err)
			datas := make(map[string]interface{}, 0)
			datas["error"] = library.Strval(err)
			c.Data["json"] = datas
			c.ServeJSON()
		}
	}()

	input := c.getFilterInput()
	if !input.IsValid("id") {
		messageInfo := c.GetValidationMessage(input)
		panic(errors.New(messageInfo))
	}
	// 在进行删除处理之前进行检查
	c.validate4Delete(input, nil)

	// delete
	c.remove(input, nil)

	c.MakeJsonResult("", "0", "", nil)

}

func (c *FormController) GetValidationMessage(input models.Input) string {
	messageInfo := ""
	errors := input.GetMessages()
	for _, err := range errors {
		messageInfo += fmt.Sprintf("Validation failure '%s': %s<br/>", err.Key, err.Message)
	}
	return messageInfo
}

func (c *FormController) GetTrueOrFalseDatas() []map[string]interface{} {
	trueMap := map[string]interface{}{"name": "是", "value": "1"}
	falseMap := map[string]interface{}{"name": "否", "value": "0"}
	ret := []map[string]interface{}{trueMap, falseMap}
	return ret
}

func (c *FormController) SysMsg(msg_detail string, msg_type int64, links []map[string]interface{}, auto_redirect bool) {
	c.Layout = "index.phtml"
	c.TplName = "error/message.phtml"
	if len(links) == 0 {
		links = make([]map[string]interface{}, 0)
		map1 := make(map[string]interface{}, 0)
		map1["text"] = "返回上一页"
		map1["href"] = "javascript:history.go(-1)"
		links = append(links, map1)
	}
	c.Data["ur_here"] = "系统信息"
	c.Data["msg_detail"] = msg_detail
	c.Data["msg_type"] = msg_type
	c.Data["links"] = links
	c.Data["default_url"] = links[0]["href"]
	c.Data["auto_redirect"] = auto_redirect

	return

}

//创建一个JSON格式的数据
func (c *FormController) MakeJsonResult(content string, errorCode string, message string, appendData map[string]interface{}) {
	datas := make(map[string]interface{}, 0)
	datas["error"] = errorCode
	datas["message"] = message
	datas["content"] = content
	for key, val := range appendData {
		datas[key] = val
	}
	c.Data["json"] = datas
	c.ServeJSON()
	return
}

func (c *FormController) CreateSchema(field string, name string, data map[string]interface{}, validation map[string]interface{}, form map[string]interface{}, list map[string]interface{}, search map[string]interface{}, export map[string]interface{}) map[string]interface{} {
	schema := make(map[string]interface{})
	schema["field"] = field
	schema["name"] = name
	schema["data"] = data
	schema["validation"] = validation
	schema["form"] = form
	schema["list"] = list
	schema["search"] = search
	schema["export"] = export

	c.Schema_keys = append(c.Schema_keys, field)
	return schema
}

func (c *FormController) getName() string {
	return c.FormName
}

func (c *FormController) getModel() system.IBase {
	return c.Model1
}

func (c *FormController) getSchemas() map[string]map[string]interface{} {

	c.Schema_keys = make([]string, 0)
	schemas := make(map[string]map[string]interface{})

	schemas["_id"] = c.CreateSchema("_id", "ID",
		map[string]interface{}{"type": "string", "length": "24"},
		map[string]interface{}{"required": false},
		map[string]interface{}{"input_type": "hidden", "is_show": true},
		map[string]interface{}{"is_show": true},
		map[string]interface{}{"is_show": true, "placeholder": "ID..."},
		map[string]interface{}{"is_show": true})

	schemas["__CREATE_TIME__"] = c.CreateSchema("__CREATE_TIME__", "创建时间",
		map[string]interface{}{"type": "datetime", "length": "19", "defaultValue": library.GetCurrentTime()},
		map[string]interface{}{"required": false},
		map[string]interface{}{"input_type": "datetimepicker", "is_show": false},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false})

	schemas["__MODIFY_TIME__"] = c.CreateSchema("__MODIFY_TIME__", "修改时间",
		map[string]interface{}{"type": "datetime", "length": "19", "defaultValue": library.GetCurrentTime()},
		map[string]interface{}{"required": false},
		map[string]interface{}{"input_type": "datetimepicker", "is_show": false},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false})

	schemas["__REMOVED__"] = c.CreateSchema("__REMOVED__", "__REMOVED__",
		map[string]interface{}{"type": "boolean", "length": "1", "defaultValue": false},
		map[string]interface{}{"required": false},
		map[string]interface{}{"input_type": "radio", "is_show": false, "items": c.GetTrueOrFalseDatas},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false},
		map[string]interface{}{"is_show": false})

	if c.GetSchemas != nil {
		schemas = c.GetSchemas(schemas)
	}

	return schemas
}

func (c *FormController) getPartials4List() []string {
	return []string{}
}

func (c *FormController) getListFilterInput() models.Input {

	draw := library.Trim(c.GetString("draw", "1"))
	//fmt.Println("draw:", draw)

	sort_order := library.Trim(c.Input().Get("order[0][dir]")) //asc
	if library.Empty(sort_order) {
		sort_order = "asc"
	}
	sort_column := library.Trim(c.Input().Get("order[0][column]")) // 0
	if library.Empty(sort_column) {
		sort_column = "0"
	}
	sort_by := library.Trim(c.Input().Get(fmt.Sprintf("columns[%s][name]", sort_column))) // _id
	if library.Empty(sort_by) {
		sort_by = "_id"
	}
	//fmt.Println("sort_order:", sort_order, "sort_column:", sort_column, "sort_by:", sort_by)

	start, _ := c.GetInt("start", 0)
	length, _ := c.GetInt("length", 10)

	input := models.NewInput()
	input.AddKeys("draw", draw)
	input.AddKeys("page", start/length+1)
	input.AddKeys("page_size", length)
	input.AddKeys("sort_by", sort_by)
	input.AddKeys("sort_order", sort_order)

	schemas := c.getSchemas()
	for key, field := range schemas {

		searchInfo, ok1 := field["search"].(map[string]interface{})
		dataInfo, ok2 := field["data"].(map[string]interface{})
		//fmt.Println("key:", key, "value:", field["data"], "search:", field["search"])

		if !ok1 || !ok2 {
			continue
		}

		is_show := library.Boolval(searchInfo["is_show"])
		dataType := library.Strval(dataInfo["type"])
		if !is_show {
			continue
		}

		input.AddSchema(key, field)

		if dataType == "integer" {
			keyValue := library.Trim(c.GetString(key, "0"))
			fmt.Println("getListFilterInput key:", key, "value:", keyValue)
			input.AddKeys(key, keyValue)
		} else if dataType == "array" {
			keyValue := c.GetStrings(key + "[]")
			fmt.Println("getListFilterInput key:", key, "value:", keyValue)
			input.AddKeys(key, keyValue)
		} else {
			keyValue := library.Trim(c.GetString(key, ""))
			fmt.Println("getListFilterInput key:", key, "value:", keyValue)
			input.AddKeys(key, keyValue)
		}
	}

	input.AddKeys("isValid", func(fieldName ...string) bool { return true })

	//input.AddKeys("getMessages", func() []*validation.Error { return nil })

	return input
}

func (c *FormController) getList(input models.Input) map[string]interface{} {

	list := c.getModel().GetList(input)
	c.Data["list"] = list["data"]
	c.Data["filter"] = list["filter"]
	c.Data["record_count"] = list["record_count"]
	c.Data["page_count"] = list["page_count"]

	return list
}

func (c *FormController) getList4Show(input models.Input, list map[string]interface{}) map[string]interface{} {

	if c.GetList4Show != nil {
		list = c.GetList4Show(input, list)
	}
	return list
}

func (c *FormController) getFilterInput() models.Input {

	input := models.NewInput()
	id := library.Trim(c.GetString("id", ""))
	input.AddKeys("id", id)
	fmt.Println("getFilterInput key:", "id", "value:", id)

	schemas := c.getSchemas()
	for key, field := range schemas {

		formInfo, ok1 := field["form"].(map[string]interface{})
		dataInfo, ok2 := field["data"].(map[string]interface{})

		if !ok1 || !ok2 {
			continue
		}

		is_show := library.Boolval(formInfo["is_show"])
		dataType := library.Strval(dataInfo["type"])
		if !is_show {
			continue
		}

		input.AddSchema(key, field)

		// 文件的话,专门处理
		if dataType == "file" {
			// 存在的话
			input.DelKeys(key)

			f, h, err := c.GetFile(key) //获取上传的文件
			defer f.Close()             //关闭上传的文件，不然的话会出现临时文件不能清除的情况
			if err != nil {
				panic(err)
			}

			path := ""
			fileInfo, ok3 := dataInfo["file"].(map[string]interface{})
			if ok3 {
				if filePath, ok4 := fileInfo["path"]; ok4 {
					path = strings.TrimRight(library.Strval(filePath), "/") + "/"
				}
			}
			uploadPath := fmt.Sprintf("%spublic/upload/%s", library.APP_PATH, path)
			destination := fmt.Sprintf("%s%s", uploadPath, h.Filename) //文件目录

			c.SaveToFile(key, destination) //存文件

			input.AddKeys(key, h.Filename)

		} else if dataType == "integer" {
			keyValue, _ := c.GetInt(key, 0)
			fmt.Println("getFilterInput key:", key, "value:", keyValue)
			input.AddKeys(key, keyValue)
		} else if dataType == "boolean" {
			keyValue, _ := c.GetBool(key, false)
			fmt.Println("getFilterInput key:", key, "value:", keyValue)
			input.AddKeys(key, keyValue)
		} else if dataType == "array" {
			keyValue := c.GetStrings(key + "[]")
			//fmt.Println("getFilterInput key:", key, "keyValue:", keyValue, "Input:", c.Input().Get(key+"[]"))
			input.AddKeys(key, keyValue)
		} else if dataType == "json" {
			keyValue := library.Trim(c.GetString(key, ""))
			fmt.Println("getFilterInput key:", key, "value:", keyValue)
			input.AddKeys(key, keyValue)
		} else if dataType == "html" {
			keyValue := library.Trim(c.GetString(key, ""))
			fmt.Println("getFilterInput key:", key, "value:", keyValue)
			input.AddKeys(key, keyValue)
		} else {
			keyValue := library.Trim(c.GetString(key, ""))
			fmt.Println("getFilterInput key:", key, "value:", keyValue)
			input.AddKeys(key, keyValue)
		}

	}

	input.AddKeys("isValid", func(fieldName ...string) bool {

		valid := validation.Validation{}

		for key, field := range schemas {

			formInfo, ok1 := field["form"].(map[string]interface{})
			validationInfo, ok2 := field["validation"].(map[string]interface{})

			if !ok1 || !ok2 {
				continue
			}

			is_show := library.Boolval(formInfo["is_show"])
			required := library.Boolval(validationInfo["required"])
			if !is_show {
				continue
			}

			if required {
				valid.Required(input.GetKeys(key), fmt.Sprintf("The %s is required", key))
			}

		}

		if valid.HasErrors() {
			// 如果有错误信息，证明验证没通过
			// 打印错误信息
			if len(fieldName) > 0 {
				errors := make([]*validation.Error, 0)
				for _, field := range fieldName {
					err, ok := valid.ErrorsMap[field]
					if ok {
						errors = append(errors, err)
					}
				}
				return !(len(errors) > 0)
			} else {
				input.AddKeys("getMessages", valid.Errors)
				return false
			}
		} else {
			return true
		}
	})

	//input.AddKeys("getMessages", func() []*validation.Error { return nil })

	return input
}

func (c *FormController) validate4Insert(input models.Input, row map[string]interface{}) bool {
	// do other validation
	var ret bool
	if c.Validate4Insert != nil {
		ret = c.Validate4Insert(input, row)
	}
	return ret
}

func (c *FormController) validate4Update(input models.Input, row map[string]interface{}) bool {
	// do other validation
	var ret bool
	if c.Validate4Update != nil {
		ret = c.Validate4Update(input, row)
	}
	return ret
}

func (c *FormController) validate4Delete(input models.Input, row map[string]interface{}) bool {
	// do other validation
	var ret bool
	if c.Validate4Delete != nil {
		ret = c.Validate4Delete(input, row)
	}
	return ret
}

func (c *FormController) insert(input models.Input, row map[string]interface{}) {
	c.processInsertOrUpdate(input, row, c.Ctx.Input)
	return
}

func (c *FormController) update(input models.Input, row map[string]interface{}) {
	c.processInsertOrUpdate(input, row, c.Ctx.Input)
	return
}

func (c *FormController) remove(input models.Input, row map[string]interface{}) {
	c.processRemove(input, row, c.Ctx.Input)
	return
}

//新增或更新操作
func (c *FormController) processInsertOrUpdate(input models.Input, row map[string]interface{}, beegoInput *context.BeegoInput) {
	value := beegoInput.Session("processInsertOrUpdate")
	if value == nil {
		beegoInput.CruSession.Set("processInsertOrUpdate", library.GetCurrentTime())

		if _, ok := row["_id"]; ok {
			c.getModel().UpdateRow(input, row)
		} else {
			c.getModel().InsertRow(input, row)
		}
		beegoInput.CruSession.Delete("processInsertOrUpdate")

	} else {
		panic(errors.New("上次的操作未完成，请稍候再试"))
	}
}

//删除操作
func (c *FormController) processRemove(input models.Input, row map[string]interface{}, beegoInput *context.BeegoInput) {
	c.getModel().RemoveRow(input, row)
}

// 导出操作
func (c *FormController) export(dataList []map[string]interface{}) map[string][][]string {
	excel := make(map[string][][]string, 0)
	schemas := c.getSchemas()
	fields := make([]string, 0)
	if len(schemas) > 0 {
		for _, key := range c.Schema_keys {
			field := schemas[key]
			if exportInfo, ok := field["export"].(map[string]interface{}); ok {
				if value, ok1 := exportInfo["is_show"]; ok1 {
					v := library.Boolval(value)
					if v {
						fields = append(fields, library.Strval(field["name"]))
					}
				}
			}

		}
	}
	if library.Empty(fields) {
		panic(errors.New("请设置导出的字段"))
	}
	excel["title"] = [][]string{fields}
	datas := [][]string{}

	for _, data := range dataList {
		item := []string{}
		for _, key := range c.Schema_keys {
			field := schemas[key]

			exportInfo, ok1 := field["export"].(map[string]interface{})
			dataInfo, ok2 := field["data"].(map[string]interface{})

			if !ok1 || !ok2 {
				continue
			}

			is_show := library.Boolval(exportInfo["is_show"])
			dataType := library.Strval(dataInfo["type"])
			if !is_show {
				continue
			}

			if dataType == "array" {

			} else if dataType == "json" {

			} else {

			}
			item = append(item, library.Strval(data[key]))
		}
		datas = append(datas, item)
	}
	excel["result"] = datas
	return excel
}

//protected function export(array $dataList)
//    {
//        $excel = array();

//        $fields = array();
//        $schemas = $this->getSchemas();
//        foreach ($schemas as $key => $field) {
//            if (empty($field['export']['is_show'])) {
//                continue;
//            }
//            $fields[] = $field['name'];
//        }
//        $excel['title'] = array_values($fields);
//        if (empty($excel['title'])) {
//            die('请设置导出的字段');
//        }
//        $datas = array();
//        foreach ($dataList as $data) {
//            $item = array();
//            foreach ($schemas as $key => $field) {
//                if (empty($field['export']) || empty($field['export']['is_show'])) {
//                    continue;
//                }
//                if ($field['data']['type'] == 'datetime') {
//                    $item[] = date("Y-m-d H:i:s", $data[$key]->sec);
//                } elseif ($field['data']['type'] == 'json') {
//                    if (! empty($data[$key])) {
//                        $values = array();
//                        if (! empty($field['export']) && ! empty($field['export']['fields'])) {
//                            foreach ($field['export']['fields'] as $f) {
//                                $values[] = isset($data[$key][$f]) ? $data[$key][$f] : "";
//                            }
//                        }
//                        if (! empty($values)) {
//                            $item[] = implode(" ", $values);
//                        } else {
//                            $item[] = "";
//                        }
//                    } else {
//                        $item[] = "";
//                    }
//                } else {
//                    $item[] = $data[$key];
//                }
//            }
//            $datas[] = $item;
//        }
//        $excel['result'] = $datas;

//        $fileName = date('YmdHis');
//        arrayToCVS($fileName, $excel);

//    }

func (c *FormController) buildSearch(schemas map[string]map[string]interface{}) {
	is_show := false
	for _, field := range schemas {
		if searchInfo, ok := field["search"].(map[string]interface{}); ok {
			if value, ok1 := searchInfo["is_show"]; ok1 {
				v := library.Boolval(value)
				if v {
					is_show = true
				}
			}
		}

	}
	map1 := make(map[string]interface{})
	map1["is_show"] = is_show
	c.Data["search"] = map1
}

func (c *FormController) buildList(schemas map[string]map[string]interface{}) {
	idx := 0
	orderIdx := 0
	orderBy := "desc"
	isFirst := false

	defaultSort, _ := c.Data["defaultSort"].(map[string]int)

	for key, field := range schemas {
		if listInfo, ok := field["list"].(map[string]interface{}); ok {
			if _, ok1 := listInfo["is_show"]; ok1 {
				if sortValue, ok2 := defaultSort[key]; !isFirst && ok2 {
					orderIdx = idx
					if sortValue == -1 {
						orderBy = "desc"
					} else {
						orderBy = "asc"
					}
					isFirst = true
				}
				idx++

			}
		}

	}

	map1 := make(map[string]interface{})
	map1["orderBy"] = orderBy
	map1["orderIdx"] = orderIdx
	map1["idx"] = idx
	c.Data["list"] = map1
}
