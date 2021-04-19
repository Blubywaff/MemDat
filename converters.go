package main

import "reflect"

// INTERNAL
// Handles structs for the interface converter
// TODO ensure valid
func convertStruct(data interface{}) (document, Result) {
	document := document{}

	ref := reflect.ValueOf(data)

	switch ref.Kind() {
	case reflect.Struct:
		break
	default:
		return nil, *newResult("Cannot parse data", FAILURE)
	}

	for i := 0; i < ref.NumField(); i++ {
		val := ref.Field(i)
		typeVal := ref.Type().Field(i)
		tag := typeVal.Tag

		switch val.Kind() {
		case reflect.Struct:
			doc := convertStruct(val.Interface())
			tagName := tag.Get("memdat")
			if tagName != "" {
				document[tagName] = doc
				continue
			}
			document[typeVal.Name] = doc
			continue
		case reflect.Slice:
			doc := convertSlice(val.Interface())
			tagName := tag.Get("memdat")
			if tagName != "" {
				document[tagName] = doc
				continue
			}
			document[typeVal.Name] = doc
			continue
		default:
			doc := convertPrimitive(val.Interface())
			tagName := tag.Get("memdat")
			if tagName != "" {
				document[tagName] = doc
				continue
			}
			document[typeVal.Name] = doc
			continue
		}

	}

	return document
}

// INTERNAL
// For use with Convert Struct to convert slices
func convertSlice(data interface{}) []interface{} {
	var items []interface{}

	for i := 0; i < reflect.ValueOf(data).Len(); i++ {
		var item interface{}
		if reflect.ValueOf(data).Kind() == reflect.Struct {
			item = convertStruct(reflect.ValueOf(data).Index(i).Interface())
		} else {
			item = convertPrimitive(reflect.ValueOf(data).Index(i).Interface())
		}

		items = append(items, item)
	}

	return items
}

// INTERNAL
// For use with converters to convert int, float, double, bool, string
func convertPrimitive(data interface{}) interface{} {
	return data
}
