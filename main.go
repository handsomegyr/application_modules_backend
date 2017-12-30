package main

import (
	"application_modules_backend/library"
	_ "application_modules_backend/routers"
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strings"

	//"io/ioutil"

	"github.com/astaxie/beego"
)

func main() {

	beego.AddTemplateExt("phtml")

	beego.AddFuncMap("hi", func(in string) (out string) {
		out = in + "world"
		return
	})

	beego.AddFuncMap("eq_string", func(inputtype interface{}, cmptype interface{}) (out bool) {
		type1 := library.Strval(inputtype)
		out = (type1 == library.Strval(cmptype))
		return
	})

	beego.AddFuncMap("get_filepath", func(fieldData map[string]interface{}, datatype string) string {
		path := ""
		info, ok := fieldData[datatype]
		if ok {
			fileInfo, ok1 := info.(map[string]interface{})
			if ok1 {
				p, ok2 := fileInfo["path"]
				if ok2 {
					path = strings.TrimRight(library.Strval(p), "/") + "/"
				}
			}
		}
		return path
	})

	beego.AddFuncMap("partial", func(viewPath string, tplFile string, data interface{}) string {

		//		dat, err := ioutil.ReadFile(viewPath + "/" + tplFile)
		//		if err != nil {
		//			panic(err)
		//		}
		//		return string(dat)
		var buf bytes.Buffer
		t, err := template.ParseFiles(viewPath + "/" + tplFile)
		//fmt.Println("data:", data)
		if err != nil {
			panic(err)
		}
		err = t.Execute(&buf, data)
		if err != nil {
			panic(err)
		}
		text := buf.String()
		//fmt.Println("partial:", text)
		return text
		//text = strings.Replace(text, "&nbsp;", " ", -1)
		//text = strings.Replace(text, "&rdquo;", "”", -1)
		//text = strings.Replace(text, "&ldquo;", "“", -1)
		//text = strings.Replace(text, "&quot;", "\"", -1)
		//text = strings.Replace(text, "&#39;", "'", -1)
		//text = strings.Replace(text, "&gt;", ">", -1)
		//text = strings.Replace(text, "&lt;", "<", -1)
		//text = strings.Replace(text, "&amp;", "&", -1) // Must be done last!

		//return strings.TrimSpace(text)

		//return buf.String()
	})

	beego.AddFuncMap("is_empty4map", func(inputtype map[string]interface{}, key string) (out bool) {
		val, out := inputtype[key]
		if out {
			out = library.Empty(val)
		}
		return
	})

	beego.AddFuncMap("get_val4textarea", func(datatype string, row map[string]interface{}, key string) string {
		val, ok := row[key]
		s := ""
		if datatype == "json" {
			if !ok {
				s = "{}"
			} else {
				s = library.Json_encode(val)
			}
		} else if datatype == "array" {
			//$value=empty($this->view->row[$key])?"":implode(',',$this->view->row[$key]);
			if !ok {
				s = ""
			}
		} else {
			s = library.Nl2br(library.Strval(val))
		}
		return s
	})

	beego.AddFuncMap("in_array", func(v interface{}, slice1 interface{}) bool {
		if slice1 == nil {
			return false
		}

		itemsType := reflect.TypeOf(slice1)
		itemsKind := itemsType.Kind()
		fmt.Println("type:", itemsType.String(), "kind:", itemsKind)
		if itemsKind == reflect.Slice {
			if itemsType.String() == "[]string" {
				items, _ := slice1.([]string)
				for _, elem := range items {
					if v == elem {
						return true
					}
				}
			}
		}
		return false
	})

	beego.AddFuncMap("create_map", func(m map[string]interface{}) map[string]interface{} {

		//	if(empty($field['form']['cascade'])){
		//	    $items = is_callable($field['form']['items'])?$field['form']['items']():$field['form']['items'];
		//	}else{
		//	    $cascade=$field['form']['cascade'];
		//	    $items = is_callable($field['form']['items'])?$field['form']['items']($this->view->row[$cascade]):$field['form']['items'];
		//	}

		itemsType := reflect.TypeOf(m["items"])
		fv := reflect.ValueOf(m["items"])
		itemsKind := itemsType.Kind()
		//fmt.Println("type:", itemsType, "kind:", itemsKind)
		if itemsKind == reflect.Func {

			//		if _, ok := m["cascade"]; ok {
			//			//$items = is_callable(m['items'])?m['items']():m['items'];
			//		} else {
			//			//$items = is_callable(m['items'])?m['items']():m['items'];
			//		}
			rs := fv.Call(nil)
			type1 := reflect.TypeOf(rs[0].Interface())
			if type1.Kind() == reflect.Slice {
				items, _ := rs[0].Interface().([]map[string]interface{})
				ret := make(map[string]interface{})

				for _, item := range items {
					//"name": "是", "value": "1"
					ret[library.Strval(item["value"])] = item["name"]
				}
				return ret
			} else if itemsKind == reflect.Map {
				items, _ := rs[0].Interface().(map[string]interface{})
				return items
			} else {
				ret := make(map[string]interface{}, 0)
				return ret
			}

		} else if itemsKind == reflect.Map {
			items, _ := m["items"].(map[string]interface{})
			return items
		} else if itemsKind == reflect.Slice {
			items, _ := m["items"].([]map[string]interface{})
			ret := make(map[string]interface{})

			for _, item := range items {
				//"name": "是", "value": "1"
				ret[library.Strval(item["value"])] = item["name"]
			}
			return ret
		} else {
			panic(500)
		}

	})

	beego.AddFuncMap("create_map2", func(m map[string]interface{}, row map[string]interface{}, key string) map[string]interface{} {

		//		if(empty($field['form']['cascade'])){
		//		   $items = is_callable($field['form']['items'])?$field['form']['items']($this->view->row[$key]):$field['form']['items'];
		//		}else{
		//		   $cascade=$field['form']['cascade'];
		//		   //die('$cascade'.$cascade.$this->view->row[$cascade]);
		//		   $items = is_callable($field['form']['items'])?$field['form']['items']($this->view->row[$cascade]):$field['form']['items'];
		//		}
		fv := reflect.ValueOf(m["items"])
		itemsType := reflect.TypeOf(m["items"])
		itemsKind := itemsType.Kind()
		//fmt.Println("type:", itemsType, "kind:", itemsKind)
		if itemsKind == reflect.Func {

			//		if _, ok := m["cascade"]; ok {
			//			//$items = is_callable(m['items'])?m['items']():m['items'];
			//		} else {
			//			//$items = is_callable(m['items'])?m['items']():m['items'];
			//		}
			rs := fv.Call(nil)
			type1 := reflect.TypeOf(rs[0].Interface())
			if type1.Kind() == reflect.Slice {
				items, _ := rs[0].Interface().([]map[string]interface{})
				ret := make(map[string]interface{})

				for _, item := range items {
					//"name": "是", "value": "1"
					ret[library.Strval(item["value"])] = item["name"]
				}
				return ret
			} else if itemsKind == reflect.Map {
				items, _ := rs[0].Interface().(map[string]interface{})
				return items
			} else {
				ret := make(map[string]interface{}, 0)
				return ret
			}
		} else if itemsKind == reflect.Map {
			items, _ := m["items"].(map[string]interface{})
			return items
		} else if itemsKind == reflect.Slice {
			items, _ := m["items"].([]map[string]interface{})
			ret := make(map[string]interface{})

			for _, item := range items {
				//"name": "是", "value": "1"
				ret[library.Strval(item["value"])] = item["name"]
			}
			return ret
		} else {
			panic(500)
		}

	})

	beego.AddFuncMap("create_map3", func(m map[string]map[string]interface{}, key string) map[string]interface{} {
		//fmt.Println("m-key:", m[key], "key:", key)
		if items1, ok := m[key]; ok {
			return items1
		} else {
			return nil
		}
	})

	beego.AddFuncMap("create_map4", func(m map[string]interface{}, key string) map[string]interface{} {
		//fmt.Println("m-key:", m[key], "key:", key)
		if items1, ok := m[key]; ok {
			return items1.(map[string]interface{})
		} else {
			return nil
		}
	})

	beego.AddFuncMap("create_slice", func(m map[string][]string, key string) []string {
		//fmt.Println("m-slice-key:", m[key], "key:", key)
		if items1, ok := m[key]; ok {
			return items1
		} else {
			return nil
		}
	})

	beego.AddFuncMap("url_get", func(base string, path string) (out string) {
		out = base + path
		return
	})

	beego.AddFuncMap("slice_get", func(in []string, idx int) (out string) {
		fmt.Println("index", idx, in)
		out = in[idx]
		return
	})

	beego.AddFuncMap("icons1_get", func(idx int) (out string) {
		icons1 := make([]string, 0)
		icons1 = append(icons1, "icon-cogs", "icon-bookmark-empty", "icon-briefcase", "icon-table", "icon-sitemap", "icon-user")
		out = icons1[idx%6]
		return
	})

	beego.AddFuncMap("icons2_get", func(idx int) (out string) {
		icons2 := make([]string, 0)
		icons2 = append(icons2, "icon-comments", "icon-coffee", "icon-time", "icon-envelope-alt", "icon-group", "icon-user")
		out = icons2[idx%6]
		return
	})

	beego.AddFuncMap("inc_index", func(idx *int, step int) (out string) {
		*idx = *idx + step
		return ""
	})

	beego.AddFuncMap("inc_index", func(idx *int, step int) (out string) {
		*idx = *idx + step
		return ""
	})

	beego.Run()
}
