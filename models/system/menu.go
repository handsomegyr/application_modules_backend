package system

import (
	"application_modules_backend/library"
	"application_modules_backend/models"
	//"fmt"
	"strings"
)

type Menu struct {
	models.Base
	Sort_keys map[string][]string
}

func NewMenu() *Menu {
	u := &Menu{}
	u.TableName = "menu"
	u.Sort_keys = make(map[string][]string, 0)
	return u
}

func (t *Menu) GetDefaultSort() []map[string]interface{} {
	sort1 := make(map[string]interface{})
	sort1["show_order"] = 1

	sort2 := make(map[string]interface{})
	sort2["_id"] = -1

	return []map[string]interface{}{sort1, sort2}
}

func (t *Menu) getQuery() []map[string]interface{} {
	query1 := make(map[string]interface{})
	query1["is_show"] = true
	return []map[string]interface{}{query1}
}

func (t *Menu) GetPrivilege(menu_list []string, requestUrl string) map[string]map[string]interface{} {
	priv_arr := make(map[string]map[string]interface{})
	sort1 := t.GetDefaultSort()
	fields1 := make(map[string]interface{})
	/* 获取权限的分组数据 */
	query1 := make(map[string]interface{})
	query1["pid"] = ""
	res := t.FindAll([]map[string]interface{}{query1}, sort1, fields1)
	t.Sort_keys = make(map[string][]string, 0)

	keys := make([]string, 0, len(res))
	for idx, rows := range res {
		key := library.Strval(rows["_id"])
		rows["relevance"] = ""
		is_active := false
		if !library.Empty(requestUrl) && !library.Empty(rows["url"]) && strings.Contains(library.Strval(rows["url"]), requestUrl) {
			//fmt.Println("rows url :", rows["url"], "requestUrl:", requestUrl)
			is_active = true
		}
		item := make(map[string]interface{})
		item["name"] = rows["name"]
		item["relevance"] = rows["relevance"]
		item["url"] = rows["url"]
		item["is_show"] = rows["is_show"]
		item["priv"] = make(map[string]interface{})
		item["is_active"] = is_active
		item["icon"] = idx
		priv_arr[key] = item
		keys = append(keys, key)
	}
	//fmt.Println("menu pid=0 :", keys)
	t.Sort_keys["pid"] = keys

	/* 按权限组查询底级的权限名称 */
	queryIn := make(map[string]interface{})
	queryIn["$in"] = library.Array_keys(priv_arr)
	query2 := make(map[string]interface{})
	query2["pid"] = queryIn
	result := t.FindAll([]map[string]interface{}{query2}, sort1, fields1)

	keys4SamePid := make(map[string][]string, 0)
	for _, priv := range result {
		key := library.Strval(priv["_id"])
		priv["relevance"] = ""
		is_active := false
		if !library.Empty(requestUrl) && !library.Empty(priv["url"]) && strings.Contains(library.Strval(priv["url"]), requestUrl) {
			//fmt.Println("priv url :", priv["url"], "requestUrl:", requestUrl)
			is_active = true
		}

		item := make(map[string]interface{})
		item["name"] = priv["name"]
		item["relevance"] = priv["relevance"]
		item["url"] = priv["url"]
		item["is_show"] = priv["is_show"]
		item["priv"] = make(map[string]interface{})
		item["is_active"] = is_active
		item["icon"] = 0

		pid := library.Strval(priv["pid"])
		keys4SamePid[pid] = append(keys4SamePid[pid], key)

		item3 := priv_arr[pid]
		item2, _ := item3["priv"].(map[string]interface{})
		item["icon"] = len(item2)
		item2[key] = item
		item3["priv"] = item2
		if is_active {
			item3["is_active"] = is_active
		}
		priv_arr[pid] = item3

	}

	for key, val := range keys4SamePid {
		t.Sort_keys[key] = val
	}

	// 将同一组的权限使用 "," 连接起来，供JS全选
	j := 0
	for action_id, action_group := range priv_arr {
		i := 0
		priv, _ := action_group["priv"].(map[string]interface{})

		item := action_group
		item["priv_list"] = strings.Join(library.Array_keys(priv), ",")
		idx := 0
		newPriv := make(map[string]interface{})
		for key, privItem := range priv {
			cando := false
			if library.In_array(key, menu_list) {
				cando = true
			}
			item2, _ := privItem.(map[string]interface{})
			item2["cando"] = cando
			item2["icon"] = idx
			idx = idx + 1
			if cando {
				i++
			}
			newPriv[key] = item2
		}
		if j == 0 {
			//fmt.Println(newPriv)
			j++
		}

		item["priv"] = newPriv

		if i > 0 {
			item["cando"] = true
		} else {
			item["cando"] = false
		}

		priv_arr[action_id] = item
	}

	return priv_arr

	//	// 排序处理
	//	priv_arr_slice := make([]map[string]interface{}, 0)
	//	for _, key := range keys {
	//		temp, _ := priv_arr[key].(map[string]interface{})
	//		priv, _ := temp["priv"].(map[string]interface{})

	//		sub_priv_arr_slice := make([]map[string]interface{}, 0)
	//		for _, sub_key := range keys4SamePid[key] {
	//			sub_temp, _ := priv[sub_key].(map[string]interface{})
	//			sub_priv_arr_slice = append(sub_priv_arr_slice, sub_temp)
	//		}
	//		temp["priv"] = sub_priv_arr_slice
	//		priv_arr_slice = append(priv_arr_slice, temp)
	//	}

	//	//ret := make(map[string]interface{})
	//	//ret["priv_arr"] = priv_arr_slice

	//	return priv_arr_slice
}
