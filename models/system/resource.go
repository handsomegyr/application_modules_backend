package system

import (
	"application_modules_backend/library"
	"application_modules_backend/models"
	"fmt"
	"sort"
	"strings"
)

type Resource struct {
	models.Base
}

func NewResource() *Resource {
	u := &Resource{}
	u.TableName = "resource"
	return u
}

func (t *Resource) GetDefaultSort() []map[string]interface{} {
	sort1 := make(map[string]interface{})
	sort1["module"] = 1

	sort2 := make(map[string]interface{})
	sort2["controller_name"] = 1

	sort3 := make(map[string]interface{})
	sort3["action_name"] = 1
	return []map[string]interface{}{sort1, sort2, sort3}
}

func (t *Resource) getQuery() []map[string]interface{} {
	query := make(map[string]interface{})
	return []map[string]interface{}{query}
}

func (t *Resource) GetPrivilege(module string, operation_list []string) map[string]map[string]interface{} {

	sort1 := t.GetDefaultSort()

	fields := make(map[string]interface{})
	query := make(map[string]interface{})
	query["module"] = module
	list := t.FindAll([]map[string]interface{}{query}, sort1, fields)

	resources := make(map[string][]map[string]interface{}, 0)
	for _, item := range list {
		// admin_form||表管理
		key := fmt.Sprintf("%s_%s||%s", item["module"], item["controller"], item["controller_name"])
		// [name] => 删除
		// [method] => remove
		// [key] => admin_form::remove
		data := make(map[string]interface{}, 0)
		data["name"] = item["action_name"]
		data["method"] = item["action"]
		data["relevance"] = ""
		data["key"] = fmt.Sprintf("%s_%s::%s", item["module"], item["controller"], data["method"])
		resources[key] = append(resources[key], data)
	}
	//fmt.Println(resources["admin_form||表管理"])

	/* 获取权限的分组数据 */
	priv_arr := make(map[string]map[string]interface{})
	for rows, _ := range resources {
		infoArr := strings.Split(rows, "||")

		data := make(map[string]interface{}, 0)
		data["name"] = infoArr[1]
		data["relevance"] = ""
		data["method"] = ""
		data["key"] = infoArr[0]
		priv_arr[rows] = data
	}

	/* 按权限组查询底级的权限名称 */
	for key, item := range resources {

		map1 := make(map[string]map[string]interface{}, 0)
		for _, priv := range item {
			data := make(map[string]interface{}, 0)
			data["name"] = priv["name"]
			data["relevance"] = priv["relevance"]
			data["method"] = priv["method"]
			data["key"] = priv["key"]
			actionKey := library.Strval(priv["key"])
			if library.In_array(actionKey, operation_list) {
				data["cando"] = true
			} else {
				data["cando"] = false
			}
			map1[actionKey] = data
		}
		priv_arr[key]["priv"] = map1
	}

	// 将同一组的权限使用 "," 连接起来，供JS全选
	for action_id, action_group := range priv_arr {
		privMap, _ := action_group["priv"].(map[string]interface{})
		priv_arr[action_id]["priv_list"] = strings.Join(library.Array_keys(privMap), ",")
		//		for key, _ := range privMap {
		//			data, _ := privMap[key].(map[string]interface{})
		//					if library.In_array(key, operation_list) {
		//						data["cando"] = true
		//					} else {
		//						data["cando"] = false
		//					}
		//			privMap[key] = data
		//		}
		//		priv_arr[action_id]["priv"] = privMap
		// 去掉错误模块
		infoArr := strings.Split(action_id, "||")
		if library.In_array(infoArr[0], []string{"admin_error", "admin_form", "admin_index"}) {
			delete(priv_arr, action_id)
		}
	}
	//fmt.Println(priv_arr)
	//ksort($priv_arr)
	//priv_arr = library.Ksort(priv_arr)

	/* 排序输出 */
	var keys []string
	for k := range priv_arr {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	rmap1 := make(map[string]map[string]interface{}, 0)
	for _, k := range keys {
		//fmt.Println("Key:", k, "Value:", m[k])
		rmap1[k] = priv_arr[k]

	}

	return rmap1
}
