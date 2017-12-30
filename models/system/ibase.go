package system

import (
	"application_modules_backend/models"
)

type IBase interface {
	// 默认排序
	GetDefaultSort() []map[string]interface{}
	// 默认查询条件
	GetQuery() []map[string]interface{}
	// 获取单条信息
	GetInfoById(id string) map[string]interface{}
	//获取列表信息
	GetList(input models.Input) map[string]interface{}
	//获取所有列表信息
	GetAllList(input models.Input) []map[string]interface{}
	//获取空行数据
	GetEmptyRow(input models.Input) map[string]interface{}
	//更新操作
	UpdateRow(input models.Input, row map[string]interface{})
	//插入操作
	InsertRow(input models.Input, row map[string]interface{})
	//删除操作
	RemoveRow(input models.Input, row map[string]interface{})
}
