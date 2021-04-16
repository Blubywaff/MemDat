package database

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
)

// INTERNAL
// used to store data in the database
// should not be used outside this pkg
type document map[string]interface{}

type indexElement struct {
	Document *document
	Value    string
}

type index struct {
	Field string
	Index []indexElement
}

type Database struct {
	documents []document
	indexes   []index
}

// INTERNAL
// Finds index of value in the index
func (i *index) findPlace(value string) int {
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
func (i *index) contains(value string) bool {
	p := i.findPlace(value)
	if p == -1 {
		return false
	}
	return i.Index[p].Value == value
}

// INTERNAL
// adds document to index
// TODO enforce unique value
func (i *index) add(document *document) bool {
	value := (*document)[i.Field].(string)
	place := int(math.Max(float64(i.findPlace(value)), 0))
	i.Index = append(append(i.Index[0:place], indexElement{document, value}), i.Index[place:]...)
	return true
}

// INTERNAL
// gets document which has the specific value
// TODO needs improvement
func (i *index) findDocument(value string) *document {
	return i.Index[i.findPlace(value)].Document
}

// INTERNAL
// Shortcut check for if an index exists
func (d *Database) hasIndex(field string) bool {
	return d.findIndexIndex(field) != -1
}

// INTERNAL
// Shortcut for searching for a document with known ObjectId
func (d *Database) findDocumentById(objectId string) *document {
	return d.findIndex("ObjectId").findDocument(objectId)
}

// INTERNAL
// Searches for document
func (d *Database) findDocument(key string, value string) *document {
	return d.findIndex(key).findDocument(value)
}

// INTERNAL
// Gets specific index based on its field
func (d *Database) findIndex(field string) *index {
	ii := d.findIndexIndex(field)
	if ii == -1 {
		return nil
	}
	return &d.indexes[ii]
}

// INTERNAL
// Gets the index (in index array) of the index that corresponds to field
func (d *Database) findIndexIndex(field string) int {
	for i, index := range d.indexes {
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
	d.indexes = append(d.indexes, index{field, make([]indexElement, len(d.documents))})
	for _, document := range d.documents {
		d.addDocumentToIndex(field, &document)
	}
	return true
}

// INTERNAL
// Adds the document to the databases indexes
func (d *Database) addDocumentToIndex(field string, document *document) bool {
	if !d.hasIndex(field) {
		return false
	}
	return d.findIndex(field).add(document)
}

// INTERNAL
// For handling adding to the database after the input has been parsed
func (d *Database) addDocument(document document) {
	d.documents = append(d.documents, document)
	for s, _ := range document {
		d.addDocumentToIndex(s, &document)
	}
}

// INTERNAL
// Generates 32 len hexadecimal ids similar to uuid
func (d *Database) generateId() string {
	id := ""
	for id == "" || d.findIndex("ObjectId").contains(id) {
		id = fmt.Sprintf("%08x%08x%08x%08x", rand.Intn(4294967296), rand.Intn(4294967296), rand.Intn(4294967296), rand.Intn(4294967296))
	}
	return id
}

// PUBLIC
// for use by others to add to the database
// may change to internal to create naming and regularity among public functions
// TODO ensure document convert succeeded
func (d *Database) Add(data interface{}) {
	document := convertStruct(data)
	document["ObjectId"] = d.generateId()

	d.addDocument(document)
}

func (d *Database) Get(data interface{}) {

}

// INTERNAL
// Handles structs for the interface converter
// TODO ensure valid
func convertStruct(data interface{}) document {
	document := document{}

	ref := reflect.ValueOf(data)

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
		items = append(items, convertPrimitive(reflect.ValueOf(data).Index(i).Interface()))
	}

	return items
}

// INTERNAL
// For use with converters to convert int, float, double, bool, string
func convertPrimitive(data interface{}) interface{} {
	return data
}

// PUBLIC
// Creates new database
func NewDatabase() *Database {
	database := Database{
		make([]document, 0),
		make([]index, 0),
	}
	database.addIndex("ObjectId")
	return &database
}
