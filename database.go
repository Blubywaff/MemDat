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
	index := d.findIndex("ObjectId")
	if index == nil {
		return nil
	}
	return index.findDocument(objectId)
}

// INTERNAL
// Searches for document
func (d *Database) findDocument(key string, value string) *document {
	index := d.findIndex(key)
	if index == nil {
		return nil
	}
	return index.findDocument(value)
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
func (d *Database) addIndex(field string) Result {
	if d.hasIndex(field) {
		return *newResult("Index Already Exists", FAILURE)
	}
	d.Indexes = append(d.Indexes, index{field, make([]indexElement, len(d.Documents))})
	for _, document := range d.Documents {
		d.addDocumentToIndex(field, &document)
	}
	return *newResult("Created Index: "+field, SUCCESS)
}

// INTERNAL
// Adds the document to the databases Indexes
func (d *Database) addDocumentToIndex(field string, document *document) Result {
	if !d.hasIndex(field) {
		return *newResult("Index of field: '"+field+"' does not exist", FAILURE)
	}
	res := d.findIndex(field).add(document)
	return res
}

// INTERNAL
// For handling adding to the database after the input has been parsed
// Should fail if the document cannot be added to one of the indexes
func (d *Database) addDocument(document document) Result {
	d.Documents = append(d.Documents, document)
	for s, _ := range document {
		if d.hasIndex(s) {
			indexResult := d.addDocumentToIndex(s, &document)
			if indexResult.IsError() {
				d.removeDocument(document)
				return *newResult("Failed to add document to Index: "+s, FAILURE)
			}
		}
	}
	return *newResult("Added document to Database", SUCCESS)
}

//INTERNAL
// Used to remove all references to a document in the database
// Should NOT fail on documents with partial addition
func (d *Database) removeDocument(document document) Result {
	dptr := d.findDocumentById(document["ObjectId"].(string))
	if dptr == nil {
		return *newResult("Could not find document: "+document["ObjectId"].(string), FAILURE)
	}
	d.removeDocumentFromIndex(dptr)
	for i := 0; i < len(d.Documents); i++ {
		if d.Documents[i]["ObjectId"].(string) == document["ObjectId"].(string) {
			d.Documents = append(d.Documents[:i], d.Documents[i+1:]...)
		}
	}
	return *newResult("Removed Document: "+document["ObjectId"].(string), SUCCESS)
}

//INTERNAL
// Used to remove all references in the indexes
func (d *Database) removeDocumentFromIndex(dptr *document) Result {
	for _, index := range d.Indexes {
		res := index.removeDocument(dptr)
		if res.IsError() {
			return *newResult("Could not remove from index\n"+res.Result(), FAILURE)
		}
	}
	return *newResult("Removed Document: "+(*dptr)["ObjectId"].(string), SUCCESS)
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

// PUBLIC
// for use by the end user to interact with the database
// TODO - make this
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
