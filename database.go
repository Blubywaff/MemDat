package main

import (
	"fmt"
	"math/rand"
)

type Database struct {
	Documents []document
	Indexes   []index
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
	return &d.Indexes[ii]
}

// INTERNAL
// Gets the index (in index array) of the index that corresponds to field
func (d *Database) findIndexIndex(field string) int {
	for i, index := range d.Indexes {
		if index.Field == field {
			return i
		}
	}
	return -1
}

// INTERNAL
// For use by the database to handle creation of new Indexes
func (d *Database) addIndex(field string) bool {
	if d.hasIndex(field) {
		return false
	}
	d.Indexes = append(d.Indexes, index{field, make([]indexElement, len(d.Documents))})
	for _, document := range d.Documents {
		d.addDocumentToIndex(field, &document)
	}
	return true
}

// INTERNAL
// Adds the document to the databases Indexes
func (d *Database) addDocumentToIndex(field string, document *document) bool {
	if !d.hasIndex(field) {
		return false
	}
	return d.findIndex(field).add(document)
}

// INTERNAL
// For handling adding to the database after the input has been parsed
func (d *Database) addDocument(document document) {
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
