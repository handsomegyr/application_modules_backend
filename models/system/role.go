package system

import (
	"application_modules_backend/library"
	"application_modules_backend/models"
	"errors"
)

type Role struct {
	models.Base
}

func NewRole() *Role {
	u := &Role{}
	u.TableName = "role"
	u.ReorganizeFields = map[string]string{"menu_list": "[]string", "operation_list": "[]string"}
	return u
}

// 默认排序
func (t *Role) GetDefaultSort() []map[string]interface{} {
	sort := make(map[string]interface{})
	sort["_id"] = -1
	return []map[string]interface{}{sort}
}

// 默认查询条件
func (t *Role) GetQuery() []map[string]interface{} {
	query := make(map[string]interface{}, 0)
	return []map[string]interface{}{query}
}

//获取列表信息
func (t *Role) GetList(input models.Input) map[string]interface{} {
	// 分页查询
	list := t.Find(input.GetQuery(), input.GetSort(), input.GetOffset(), input.GetLimit(), nil)
	/* 记录总数 */
	input.SetRecordCount(list["total"].(int64))
	filter := input.GetFilter()

	ret := make(map[string]interface{}, 0)
	ret["data"] = list["datas"]
	ret["filter"] = filter
	ret["page_count"] = filter["page_count"]
	ret["record_count"] = filter["record_count"]

	return ret
}

func (t *Role) CheckName(id string, name string) bool {
	/* 判断是否已经存在 */
	query1 := make(map[string]interface{})
	query1["name"] = name

	query := t.GetQuery()
	query = append(query, query1)

	if !library.Empty(id) {
		subquery := make(map[string]interface{})
		subquery["$ne"] = id

		query2 := make(map[string]interface{})
		query2["_id"] = subquery
		query = append(query, query2)
	}
	num := t.Count(query)
	if num > 0 {
		panic(errors.New("角色名已存在"))
	}
	return true
}

//获取所有列表信息
func (t *Role) GetAllList(input models.Input) []map[string]interface{} {
	// 分页查询
	list := t.FindAll(input.GetQuery(), input.GetSort(), nil)
	return list
}

//获取空行数据
func (t *Role) GetEmptyRow(input models.Input) map[string]interface{} {
	data := input.GetFormData(false)
	return data
}

//新增操作
func (t *Role) InsertRow(input models.Input, row map[string]interface{}) {
	data := input.GetFormData(true)
	delete(row, "_id")
	// t.SetDebug(true)
	t.Insert(data)
}

//更新操作
func (t *Role) UpdateRow(input models.Input, row map[string]interface{}) {
	data := input.GetFormData(true)
	// t.SetDebug(true)
	if id, ok := row["_id"]; ok {
		query := make(map[string]interface{}, 0)
		query["_id"] = id
		contents := make(map[string]interface{})
		contents["$set"] = data

		t.Update([]map[string]interface{}{query}, contents, nil)
	} else {
		panic(500)
	}
}

//删除操作
func (t *Role) RemoveRow(input models.Input, row map[string]interface{}) {
	query := make(map[string]interface{}, 0)
	query["_id"] = library.Strval(input.GetKeys("id"))
	t.Remove([]map[string]interface{}{query})
}
