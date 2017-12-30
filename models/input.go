package models

import (
	"application_modules_backend/library"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/astaxie/beego/validation"
)

type Input struct {
	filter  map[string]interface{}
	schemas map[string]interface{}
	keys    map[string]interface{}
}

func NewInput() Input {
	u := Input{make(map[string]interface{}, 0), make(map[string]interface{}, 0), make(map[string]interface{}, 0)}
	return u
}

func (t *Input) IsValid(keys ...string) bool {
	if func1, ok := t.keys["isValid"]; ok {
		func2, _ := func1.(func(...string) bool)
		if len(keys) == 0 {
			return func2()
		} else {
			return func2(keys...)
		}
	}
	return false
}

func (t *Input) GetMessages() []*validation.Error {
	errors, ok := t.keys["getMessages"].([]*validation.Error)
	if ok {
		return errors
	}
	return nil
}

func (t *Input) AddKeys(key string, value interface{}) {
	t.keys[key] = value
}

func (t *Input) GetKeys(key string) interface{} {
	return t.keys[key]
}

func (t *Input) DelKeys(key string) {
	delete(t.keys, key)
}

func (t *Input) GetDefaultSort() map[string]interface{} {
	sort := make(map[string]interface{})
	sort["show_order"] = 1
	sort["_id"] = -1
	return sort
}

//过滤信息
func (t *Input) GetFilter() map[string]interface{} {
	if len(t.filter) == 0 {
		schemas := t.GetSchemas()
		if len(schemas) > 0 {
			for key, _ := range schemas {
				t.filter[key] = library.Urldecode(library.Trim(library.Strval(t.keys[key])))
			}
		}

		t.filter["sort_by"] = library.Trim(library.Strval(t.keys["sort_by"]))
		t.filter["sort_order"] = library.Trim(library.Strval(t.keys["sort_order"]))

		/* 分页大小 */
		// 每页显示数量
		if page_size, ok := t.keys["page_size"]; ok {
			t.filter["page_size"] = page_size
		} else {
			t.filter["page_size"] = 10
		}

		// 当前页数
		if page, ok := t.keys["page"]; ok {
			t.filter["page"] = page
		} else {
			t.filter["page"] = 1
		}

		// offset
		page := library.Intval(t.filter["page"])
		page_size := library.Intval(t.filter["page_size"])
		t.filter["start"] = (page - 1) * page_size
	}
	return t.filter
}

func (t *Input) SetRecordCount(record_count int64) {
	/* page 总数 */
	t.filter["record_count"] = record_count
	if record_count1, ok := t.keys["record_count"]; ok {

		record_count2, _ := record_count1.(float64)
		page_size, _ := t.filter["page_size"].(float64)

		t.filter["page_count"] = math.Ceil(record_count2 / page_size)
		t.filter["page_count"] = t.filter["page_count"]
	} else {
		t.filter["page_count"] = 1
	}

	/* 边界处理 */
	page := library.Intval(t.filter["page"])
	page_count := library.Intval(t.filter["page_count"])
	if page > page_count {
		t.filter["page"] = t.filter["page_count"]
	}
}

func (t *Input) GetSchemas() map[string]interface{} {
	return t.schemas
}

func (t *Input) AddSchema(key string, field interface{}) {
	t.schemas[key] = field
}

//根据画面条件获取查询条件
func (t *Input) GetQuery() []map[string]interface{} {

	filter := t.GetFilter()
	schemas := t.GetSchemas()

	where := make(map[string]interface{}, 0)
	if len(schemas) > 0 {
		for key, field1 := range schemas {
			if value, ok := filter[key]; ok {
				if !library.Empty(value) {

					field, _ := field1.(map[string]interface{})
					if fieldData1, ok1 := field["data"]; ok1 {
						fieldData, _ := fieldData1.(map[string]interface{})
						if fieldData["type"] == "string" && key != "_id" {
							//array('$like' => '%' . $q . '%')
							condition := make(map[string]interface{}, 0)
							condition["$like"] = "%" + library.Urldecode(library.Strval(value)) + "%"
							where[key] = condition
						} else if fieldData["type"] == "datetime" {
							datetime := library.Urldecode(library.Strval(value))
							datatimeArr := library.Explode("|", datetime)
							fromTo := make(map[string]interface{}, 0)
							for idx, v := range datatimeArr {
								if len(library.Trim(v)) > 0 {
									if idx == 0 {
										fromTo["$gte"] = v
									} else {
										fromTo["$lte"] = v
									}
								}
							}
							where[key] = fromTo
						} else if fieldData["type"] == "integer" {
							num := library.Urldecode(library.Strval(value))
							numArr := library.Explode("|", num)
							fromTo := make(map[string]interface{}, 0)
							for idx, v := range numArr {
								if len(library.Trim(v)) > 0 {
									if idx == 0 {
										fromTo["$gte"] = v
									} else {
										fromTo["$lte"] = v
									}
								}
							}
							where[key] = fromTo
						} else {
							where[key] = library.Urldecode(library.Strval(value))
						}
					}
				}
			}
		}
	}
	return []map[string]interface{}{where}
}

func (t *Input) GetSort() []map[string]interface{} {
	filter := t.GetFilter()
	sort_order := library.Strval(filter["sort_order"])
	sort_by := library.Strval(filter["sort_by"])
	// 排序方式
	sort := make(map[string]interface{})
	if "desc" == strings.ToLower(sort_order) {
		sort[sort_by] = -1
	} else {
		sort[sort_by] = 1
	}
	return []map[string]interface{}{sort}
}

func (t *Input) GetOffset() int64 {
	filter := t.GetFilter()
	offset := library.Intval(filter["start"])
	return offset
}

func (t *Input) GetLimit() int64 {
	filter := t.GetFilter()
	page_size := library.Intval(filter["page_size"])
	//fmt.Println("GetLimit:", page_size)
	return page_size
}

func (t *Input) GetFormData(is_update bool) map[string]interface{} {

	schemas := t.GetSchemas()
	data := make(map[string]interface{}, 0)
	if len(schemas) > 0 {
		for key, field1 := range schemas {
			value, ok := t.keys[key]
			field, _ := field1.(map[string]interface{})
			if fieldData1, ok1 := field["data"]; ok1 {
				fieldData, _ := fieldData1.(map[string]interface{})
				if fieldData["type"] == "string" {
					if ok && is_update {
						data[key] = library.Urldecode(library.Trim(library.Strval(value)))
					} else {
						if value1, ok1 := fieldData["defaultValue"]; ok1 {
							data[key] = library.Strval(value1)
						} else {
							data[key] = ""
						}
					}
				} else if fieldData["type"] == "integer" {
					if ok && is_update {
						data[key] = library.Intval(value)
					} else {
						if value1, ok1 := fieldData["defaultValue"]; ok1 {
							data[key] = library.Intval(value1)
						} else {
							data[key] = 0
						}
					}
				} else if fieldData["type"] == "float" {
					if ok && is_update {
						data[key] = value.(float32)
					} else {
						if value1, ok1 := fieldData["defaultValue"]; ok1 {
							data[key] = value1
						} else {
							data[key] = 0.00
						}
					}
				} else if fieldData["type"] == "decimal" {
					if ok && is_update {
						data[key] = value.(float64)
					} else {
						if value1, ok1 := fieldData["defaultValue"]; ok1 {
							data[key] = value1
						} else {
							data[key] = 0.00
						}
					}
				} else if fieldData["type"] == "datetime" {
					if ok && is_update {
						data[key] = value
					} else {
						if value1, ok1 := fieldData["defaultValue"]; ok1 {
							data[key] = value1
						} else {
							data[key] = library.GetCurrentTime()
						}
					}
				} else if fieldData["type"] == "boolean" {
					if ok && is_update {
						data[key] = library.Boolval(value)
					} else {
						if value1, ok1 := fieldData["defaultValue"]; ok1 {
							data[key] = library.Boolval(value1)
						} else {
							data[key] = false
						}
					}
				} else if fieldData["type"] == "json" {
					if ok && is_update {
						data[key] = library.Json_decode(library.Urldecode(library.Trim(library.Strval(value))), true)
					} else {
						if value1, ok1 := fieldData["defaultValue"]; ok1 {
							data[key] = value1
						} else {
							data[key] = make(map[string]interface{}, 0)
						}
					}
				} else if fieldData["type"] == "array" {
					if ok && is_update {
						if library.Is_array(value) {
							data[key] = value
						} else {
							panic(errors.New(fmt.Sprintf("%s is not array  ", key)))
						}
					} else {
						if value1, ok1 := fieldData["defaultValue"]; ok1 {
							data[key] = value1
						} else {
							data[key] = make([]string, 0)
						}
					}
				} else if fieldData["type"] == "html" {
					if ok && is_update {
						data[key] = library.Trim(library.Strval(value))
					} else {
						if value1, ok1 := fieldData["defaultValue"]; ok1 {
							data[key] = value1
						} else {
							data[key] = ""
						}
					}
				} else if fieldData["type"] == "file" {
					delete(data, key)
					if ok {
						data[key] = library.Trim(library.Strval(value))
					}
				} else {
					data[key] = ""
				}
			}

		}
	}

	if is_update {
		delete(data, "_id")
	} else {
		data["_id"] = ""
	}
	return data
}
