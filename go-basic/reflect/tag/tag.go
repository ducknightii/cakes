package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Item 反射解析 tag定义
type Item struct {
	Line int    `json:"line" export:"name:行号;priority:1"`
	Name string `json:"name" export:"name:名称"`
	Msg  string `json:"msg" export:"name:消息;priority:2"`
}

func main() {
	arr := []Item{
		{
			1,
			"哈哈",
			"111",
		},
		{
			3,
			"sss",
			"222",
		},
	}
	fmt.Println(reflectDump(&arr))
}

func reflectDump(ptr interface{}) error {
	// 获取入参的类型
	reType := reflect.TypeOf(ptr)
	fmt.Printf("reType: %#v\n", reType.Name())
	fmt.Printf("reType.Kind: %#v \nreType.Elem.Kind: %#v\n", reType.Kind(), reType.Elem().Kind())
	fmt.Printf("reType.Elem.Elem.Kind: %#v\n", reType.Elem().Elem().Kind())
	// 入参类型校验
	if reType.Kind() != reflect.Ptr || reType.Elem().Kind() != reflect.Slice || reType.Elem().Elem().Kind() != reflect.Struct {
		return fmt.Errorf("must be slice struct ptr")
	}
	var exportFields []string                      // 有序的要导出的字段名
	var exportTagNameMap = make(map[string]string) // fieldName => name
	// 打印key & 打印tag
	for i := 0; i < reType.Elem().Elem().NumField(); i++ {
		fieldName := reType.Elem().Elem().Field(i).Name
		exportTag := strings.TrimSpace(reType.Elem().Elem().Field(i).Tag.Get("export"))

		fmt.Printf("Struct Key[%d]: %s\tTag: %s\n", i, fieldName, exportTag)
		if exportTag != "" {
			// 解析 exportTag
			cuts := strings.Split(exportTag, ";")
			var name string = fieldName
			var priority int = i + 1
			for _, cut := range cuts {
				items := strings.Split(cut, ":")
				if len(items) == 2 && items[1] != "" {
					switch items[0] {
					case "name":
						name = items[1]
					case "priority":
						_p, _ := strconv.Atoi(items[1])
						if _p > 0 && _p < priority {
							priority = _p
						}

					}
				}
			}
			exportTagNameMap[fieldName] = name
			exportFields = append(exportFields[:priority-1], append([]string{fieldName}, exportFields[priority-1:]...)...)
			fmt.Printf("priority:%d\texportFields:%+v\texportTagNameMap:%+v\n", priority, exportFields, exportTagNameMap)

		}

	}

	reVal := reflect.ValueOf(ptr).Elem()
	fmt.Printf("reVal:%#v, type:%T\n", reVal, reVal)

	fmt.Printf("exportFields:%+v\texportTagNameMap:%+v\n", exportFields, exportTagNameMap)
	var res [][]string = make([][]string, 0, reVal.Len())
	for i := 0; i < reVal.Len(); i++ {
		v := reVal.Index(i)
		fmt.Printf("reVal[%d]: %#v\t type:%T\n", i, v, v)

		var rowArr []string = make([]string, 0, len(exportFields))
		/*
			vType := v.Type()

			for j := 0; j < v.NumField(); j++ {
				item := v.Field(j)
				itemTypeName := vType.Field(j).Name
				fmt.Printf("reVal[%d].%d: %#v\t structField:%v\n", i, j, item, itemTypeName)

			}*/
		for _, fieldName := range exportFields {
			item := v.FieldByName(fieldName)
			fmt.Printf("reVal[%d]: %#v\t structField:%v\n", i, item, fieldName)
			rowArr = append(rowArr, fmt.Sprintf("%v", item))
		}
		res = append(res, rowArr)
	}

	fmt.Println(res)
	return nil
}
