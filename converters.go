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
			doc, cres := convertStruct(val.Interface())
			if cres.IsError() {
				return nil, *newResult("Failed on sub struct: "+cres.Result(), FAILURE)
			}
			tagName := tag.Get("memdat")
			if tagName != "" {
				document[tagName] = doc
				continue
			}
			document[typeVal.Name] = doc
			continue
		case reflect.Slice:
			doc, cres := convertSlice(val.Interface())
			if cres.IsError() {
				return nil, *newResult("Failed on slice: "+cres.Result(), FAILURE)
			}
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

	return document, *newResult("Converted Struct", SUCCESS)
}

// INTERNAL
// For use with Convert Struct to convert slices
func convertSlice(data interface{}) ([]interface{}, Result) {
	var items []interface{}

	for i := 0; i < reflect.ValueOf(data).Len(); i++ {
		var item interface{}
		if reflect.ValueOf(data).Kind() == reflect.Struct {
			var cres Result
			item, cres = convertStruct(reflect.ValueOf(data).Index(i).Interface())
			if cres.IsError() {
				return nil, *newResult("Subconvert failed on element: "+string(rune(i)), FAILURE)
			}
		} else {
			item = convertPrimitive(reflect.ValueOf(data).Index(i).Interface())
		}

		items = append(items, item)
	}

	return items, *newResult("Converted Slice", SUCCESS)
}

// INTERNAL
// For use with converters to convert int, float, double, bool, string
func convertPrimitive(data interface{}) interface{} {
	return data
}
