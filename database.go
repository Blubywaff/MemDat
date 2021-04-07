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

// PUBLIC
// for use by others to add to the database
// may change to internal to create naming and regularity among public functions
// TODO reflection for other fields
func (d *Database) add(data interface{}) {
	document := Document{"ObjectId": d.generateId()}
	//mdc, ok := data.(MemDatCompatible)
	//if ok {
	//parse := mdc.MemDatConvert()

	//}
	ref := reflect.ValueOf(&data).Elem()
	for i := 0; i < ref.NumField(); i++ {
		//val := ref.Field(i)
		typeVal := ref.Type().Field(i)
		tag := typeVal.Tag
		fmt.Println(tag.Get("memdat"))
	}
	d.addDocument(document)
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
