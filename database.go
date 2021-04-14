package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
)

type Document map[string]interface{}

type IndexElement struct {
	Document *Document
	Value    string
}

type Index struct {
	Field string
	Index []IndexElement
}

type Database struct {
	Documents []Document
	Indexes   []Index
}

type MemDatCompatible interface {
	// Converts type to json format without using reflection ideally for efficiency
	MemDatConvert() string
}

// INTERNAL
// Finds index of value in the index
func (i Index) findPlace(value string) int {
	bot := 0
	top := len(i.Index)
	cur := -1
	for bot != top {
		cur = (bot+top)/2 + bot
		if i.Index[cur].Value > value {
			top = cur
		} else {
			bot = cur
		}
	}
	return cur
}

// INTERNAL
// checks if a specific value exists
func (i Index) contains(value string) bool {
	p := i.findPlace(value)
	if p == -1 {
		return false
	}
	return i.Index[p].Value == value
}

// INTERNAL
// adds document to index
// TODO enforce unique value
func (i Index) add(document *Document) bool {
	value := (*document)[i.Field].(string)
	place := int(math.Max(float64(i.findPlace(value)), 0))
	i.Index = append(append(i.Index[0:place], IndexElement{document, value}), i.Index[place:]...)
	return true
}

// INTERNAL
// gets Document which has the specific value
func (i Index) findDocument(value string) *Document {
	return i.Index[i.findPlace(value)].Document
}

// INTERNAL
// Shortcut check for if an index exists
func (d Database) hasIndex(field string) bool {
	return d.findIndexIndex(field) != -1
}

// INTERNAL
// Shortcut for searching for a document with known ObjectId
func (d *Database) findDocumentById(objectId string) *Document {
	return d.findIndex("ObjectId").findDocument(objectId)
}

// INTERNAL
// Gets specific index based on its field
func (d *Database) findIndex(field string) *Index {
	ii := d.findIndexIndex(field)
	if ii == -1 {
		return nil
	}
	return &d.Indexes[ii]
}

// INTERNAL
// Gets the index (in Index array) of the index that corresponds to field
func (d *Database) findIndexIndex(field string) int {
	for i, index := range d.Indexes {
		if index.Field == field {
			return i
		}
	}
	return -1
}

// INTERNAL
// For use by the database to handle creation of new indexes
func (d *Database) addIndex(field string) bool {
	if d.hasIndex(field) {
		return false
	}
	d.Indexes = append(d.Indexes, Index{field, make([]IndexElement, len(d.Documents))})
	for _, document := range d.Documents {
		d.addDocumentToIndex(field, &document)
	}
	return true
}

// INTERNAL
// Adds the document to the databases indexes
func (d *Database) addDocumentToIndex(field string, document *Document) bool {
	if !d.hasIndex(field) {
		return false
	}
	fmt.Println("")
	return d.findIndex(field).add(document)
}

// INTERNAL
// For handling adding to the database after the input has been parsed
func (d *Database) addDocument(document Document) {
	d.Documents = append(d.Documents, document)
	for s, _ := range document {
		d.addDocumentToIndex(s, &document)
	}
}

// INTERNAL
// Generates 32 len hexadecimal ids similar to uuid
func (d *Database) generateId() string {
	id := ""
	for id == "" || d.findIndex("ObjectId").contains(id) {
		id = fmt.Sprintf("%x%x%x%x", rand.Intn(4294967296), rand.Intn(4294967296), rand.Intn(4294967296), rand.Intn(4294967296))
	}
	return id
}

// INTERNAL
// Removes newlines spaces and tabs from input
/*
func jsonFormat(str string) string {
	quote := false
	var removes []int
	for i, _ := range str {
		if str[i] == '"' {
			quote = !quote
			continue
		}
		if str[i] == ' ' || str[i] == '\n' || str[i] == '\t' {
			removes = append(removes, i)
		}
	}
	var b strings.Builder
	b.Grow(len(str) - len(removes))
	for i, remove := range removes {
		var s int
		e := remove
		if i != 0 {
			s = removes[i]
		}
		if i != len(removes) - 1 {
			e = len(str)
		}
		_, _ = fmt.Fprint(&b, str[s:e])
	}
	return b.String()
}
*/

// INTERNAL
// Checks if input json is valid
// TODO finish this
/*
func jsonValid(str string) bool {
	secChars := map[uint8]int{
		'{': 0,
		'}': 0,
		'[': 0,
		']': 0,
	}
	escChars := []uint8{
		'\b',
		'\f',
		'\n',
		'\r',
		'\t',
		'\\',
	}
	quote := false
	for i := 0; i < len(str); i++ {
		skip := false
		for c, _ := range secChars {
			if str[i] == c {
				secChars[c]++
				if secChars['}'] > secChars['{'] || secChars[']'] > secChars['['] {
					// More Closes than Opens
					return false
				}
				skip = true
				break
			}
		}
		if !skip {
			if str[i] == '"' {
				if quote {
					if str[i-1] == '\\' {
						continue
					}
					quote = false
					if str[i+1] != ',' && str[i+1] != ':' {
						return false
					}
				} else {
					quote = true
				}
			}
		}
		if !skip && !quote{
			for _, char := range escChars {
				if str[i] == char {
					return false
				}
			}
		}
	}
	return str != ""
}

// INTERNAL
// Prepares input jsons for conversion to database object
func jsonDeparse(str string) string {
	f := jsonFormat(str)
	if !jsonValid(f) {
		return ""
	}

}
*/

// PUBLIC
// for use by others to add to the database
// may change to internal to create naming and regularity among public functions
// TODO reflection for other fields
func (d *Database) add(data interface{}) {
	document := Document{"ObjectId": d.generateId()}
	//mdc, ok := data.(MemDatCompatible)
	/*if ok {
		parse := mdc.MemDatConvert()
		parse := jsonDeparse(parse)
	}*/

	fmt.Println(convertStruct(data))

	/*
		ref := reflect.ValueOf(&data).Elem().Elem()
		fmt.Println(ref.Kind())
		//fmt.Println(ref.Elem().Kind())
		for i := 0; i < ref.NumField(); i++ {
			val := ref.Field(i)
			switch val.Kind() {
			case reflect.Struct:

			}
			typeVal := ref.Type().Field(i)
			tag := typeVal.Tag
			fmt.Println("val:", val)
			fmt.Println("isArray:", val.Kind() == reflect.Struct)
			fmt.Println("typeVal:", typeVal.Type)
			fmt.Println("tags:", tag)
		}
	*/
	d.addDocument(document)
}

// INTERNAL
// Converts input from add into document for use by database
/*
func convertInterface(data interface{}) Document {
	fmt.Println("Convert Interface")
	document := Document{}
	ref := reflect.ValueOf(data)
	var refE reflect.Value
	fmt.Println("ref kind:", ref.Kind())
	if ref.Kind() == reflect.Struct {
		convertStruct(data)
	} else if ref.Kind() == reflect.Interface {
		refE = ref.Elem()
		fmt.Println("refE kind:", refE.Kind())
	}
	fmt.Println("Ref Kind (ac):", ref.Kind())
	for i := 0; i < ref.NumField(); i++ {
		fmt.Println("Check field:", i)
		val := ref.Field(i)
		fmt.Println("val kind:", val.Kind())
		typeVal := ref.Type().Field(i)
		tag := typeVal.Tag
		switch val.Kind() {
		case reflect.Struct:
			doc := convertStruct(reflect.Indirect(val))
			tagName := tag.Get("memdat")
			if tagName != "" {
				document[tagName] = doc
				continue
			}
			document[typeVal.Name] = doc
			continue
		}
		fmt.Println("val:", val)
		fmt.Println("isArray:", val.Kind() == reflect.Struct)
		fmt.Println("typeVal:", typeVal.Type)
		fmt.Println("tags:", tag)
	}
	fmt.Println("End Convert Struct")
	return document
}
*/

// INTERNAL
// Handles structs for the interface converter
func convertStruct(data interface{}) Document {
	fmt.Println("Convert Struct:", data)
	document := Document{}

	ref := reflect.ValueOf(data)
	fmt.Println("ref kind:", ref.Kind(), ref.Type())

	for i := 0; i < ref.NumField(); i++ {
		fmt.Println("checking:", i)

		val := ref.Field(i)
		typeVal := ref.Type().Field(i)
		tag := typeVal.Tag

		fmt.Println("Val Details:")
		fmt.Println("\tVal Interface:", val.Interface())
		fmt.Println("\tRereflect kind:", reflect.ValueOf(val.Interface()).Type())
		fmt.Println("\tVal Indirect:", reflect.Indirect(val))

		fmt.Println("Field Kind:", val.Kind())

		fmt.Println("Field Kind ap:", reflect.Indirect(val).Kind())

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

	fmt.Println("End Convert Struct")
	return document
}

// INTERNAL
// For use with Convert Struct to convert slices
func convertSlice(data interface{}) Document {
	fmt.Println("Convert Slice:", data)
	document := Document{}

	fmt.Println("End Convert Slice")
	return document
}

// INTERNAL
// For use with converters to convert int, float, double, bool, string
func convertPrimitive(data interface{}) interface{} {
	fmt.Println("Convert Primitive:", data)
	//document := Document{}

	fmt.Println(reflect.ValueOf(data).Type())

	fmt.Println("End Convert Primitive")
	return data
}

// PUBLIC
// Creates new database
func newDatabase() *Database {
	database := Database{
		make([]Document, 0),
		make([]Index, 0),
	}
	database.addIndex("ObjectId")
	return &database
}
