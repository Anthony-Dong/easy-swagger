package swagger

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	objectType = "object"
	sliceType  = "array"
	mapType    = "object"
)

func AddProperties(tp reflect.Type, defines map[string]*Definition, fieldProperty *Property) {
	if defines == nil {
		panic("the receive map is nil")
	}
	for tp.Kind() == reflect.Ptr { //消除ptr
		tp = tp.Elem()
	}
	switch tp.Kind() {
	case reflect.Struct:
		if fieldProperty != nil { // 如果fieldProperty不为空，则说明可能存在引用
			item := getEndItem(fieldProperty) // item和map不会同时存在的
			endMap := getEndMap(fieldProperty)
			if item != nil && isEqueal(item.Type, objectType) {
				item.Ref = getRef(tp)
			} else if endMap != nil && isEqueal(endMap.Type, objectType) {
				endMap.Ref = getRef(tp)
			} else {
				fieldProperty.Type = objectType
				fieldProperty.Ref = getRef(tp)
			}
		}
		_, isExist := defines[getStructName(tp)] // 如果是结构体，直接先看看有没有加载过，1、没有加载直接去load
		if isExist {
			return // 存在代表加载过了
		}
		define := new(Definition)                  // 初始化结构体
		defines[getStructName(tp)] = define        // 添加到结构体中
		define.Properties = map[string]*Property{} // 添加结构体信息
		define.Type = objectType                   // 结构体
		field := tp.NumField()                     // 遍历字段
		for x := 0; x < field; x++ {
			field := tp.Field(x)
			if isIgnoreField(field) { // 忽略字段：1、非public字段，2、json ignore
				continue
			}
			property := new(Property)
			define.Properties[getFieldName(field)] = property
			property.Description = getDesc(field)
			fieldType := field.Type
			for fieldType.Kind() == reflect.Ptr { //获取真实类型
				fieldType = fieldType.Elem()
			}
			switch fieldType.Kind() {
			case reflect.Struct: // 结构体：1、字段为结构体，直接递归
				property.Type = objectType
				AddProperties(fieldType, defines, property)
				continue
			case reflect.Interface:
				property.Type = objectType // interface  模式 就走  object，什么都接受
				continue
			case reflect.Slice:
				property.Type = sliceType
				AddProperties(fieldType, defines, property)
				continue
			case reflect.Map:
				property.Type = mapType
				AddProperties(fieldType, defines, property)
				continue
			default:
				kind, isExist := mapping[fieldType.Kind()] // 这里只处理基本类型
				if isExist {
					property.Type = kind
					property.Format = fieldType.Kind().String()
				}
				continue
			}
		}
	case reflect.Map:
		if fieldProperty == nil {
			panic("Map fieldProperty is nil")
		}
		valueType := tp.Elem()
		for valueType.Kind() == reflect.Ptr {
			valueType = valueType.Elem()
		}
		item := newMap(fieldProperty)
		switch valueType.Kind() {
		case reflect.Slice:
			item.Type = sliceType
			AddProperties(valueType, defines, fieldProperty)
			return
		case reflect.Struct:
			item.Type = objectType
			AddProperties(valueType, defines, fieldProperty)
			return
		case reflect.Map:
			item.Type = mapType
			AddProperties(valueType, defines, fieldProperty)
			return
		default:
			kind, isExist := mapping[valueType.Kind()]
			if isExist {
				item.Type = kind
				return
			}
			return
		}
	case reflect.Slice:
		if fieldProperty == nil {
			panic("Slice fieldProperty is nil")
		}
		item := newItem(fieldProperty) // 切片需要初始化item，初始化逻辑是，在末端生成一个，然后设置类型
		valueType := tp.Elem()
		for valueType.Kind() == reflect.Ptr { //获取真实类型
			valueType = valueType.Elem()
		}
		switch valueType.Kind() {
		case reflect.Slice:
			item.Type = sliceType
			AddProperties(valueType, defines, fieldProperty)
			return
		case reflect.Struct:
			item.Type = objectType
			AddProperties(valueType, defines, fieldProperty)
			return
		case reflect.Map:
			item.Type = mapType
			AddProperties(valueType, defines, fieldProperty)
			return
		default:
			kind, isExist := mapping[valueType.Kind()]
			if isExist {
				item.Type = kind
			}
			return
		}
	default:
		if fieldProperty == nil {
			panic("default fieldProperty is nil")
		}
		kind, isExist := mapping[tp.Kind()]
		if isExist {
			fieldProperty.Type = kind
			fieldProperty.Format = tp.Kind().String()
		}
		return
	}
}

/**
public 的 pkg_path == ""
*/
func isIgnoreField(field reflect.StructField) bool {
	if field.PkgPath != "" {
		return true
	}
	return strings.Compare(field.Tag.Get("json"), "-") == 0
}

func getEndItem(property *Property) *Items {
	if property == nil {
		return nil
	}
	var item = property.Items
	for item != nil {
		if item.Items == nil {
			return item
		}
		item = item.Items
	}
	return item
}

func getEndMap(property *Property) *MapValueProperty {
	if property == nil {
		return nil
	}
	var item = property.MapValueProperties
	if item != nil {
		if item.MapValueProperties == nil {
			return item
		}
		item = item.MapValueProperties
	}
	return item
}

func newItem(fieldProperty *Property) *Items {
	var item *Items
	if fieldProperty.Items == nil {
		item = new(Items)
		fieldProperty.Items = item
	} else {
		newItem := fieldProperty.Items
		item = new(Items)
		for newItem != nil {
			if newItem.Items == nil {
				newItem.Items = item
				break
			}
			newItem = newItem.Items
		}
	}
	return item
}
func newMap(fieldProperty *Property) *MapValueProperty {
	var item *MapValueProperty
	if fieldProperty.MapValueProperties == nil {
		item = new(MapValueProperty)
		fieldProperty.MapValueProperties = item
	} else {
		newItem := fieldProperty.MapValueProperties
		item = new(MapValueProperty)
		for newItem != nil {
			if newItem.MapValueProperties == nil {
				newItem.MapValueProperties = item
				break
			}
			newItem = newItem.MapValueProperties
		}
	}
	return item
}

func isEqueal(str1, str2 string) bool {
	return strings.Compare(str1, str2) == 0
}

func getStructName(tp reflect.Type) string {
	return tp.String()
}

func getDesc(field reflect.StructField) string {
	desc := field.Tag.Get("desc")
	return desc
}

func getRef(type_ reflect.Type) string {
	return fmt.Sprintf("#/definitions/%s", type_.String())
}

func GetRef(type_ reflect.Type) string {
	return getRef(type_)
}

func getFieldName(field reflect.StructField) string {
	fieldName := field.Tag.Get("json")
	if fieldName == "" {
		return field.Name
	}
	return fieldName
}
