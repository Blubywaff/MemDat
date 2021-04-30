package main

import (
	"fmt"
	"math/rand"
	"reflect"
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
// Searches for document from indexes
func (d *Database) findDocument(key string, value string) *document {
	index := d.findIndex(key)
	if index == nil {
		return nil
	}
	return index.findDocument(value)
}

// INTERNAL
// Searches for document not in index
func (d *Database) findDocumentNoIndex(key string, value interface{}) *document {
	for _, doc := range d.Documents {
		if doc[key] == value {
			return &doc
		}
	}
	return nil
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
				res := d.removeDocument(document)
				if res.IsError() {
					panic("Failed to both add and remove Document. " +
						"This SHOULD NOT EVER happen! " +
						"This will corrupt the database")
				}
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
		id = fmt.Sprintf("%08x", rand.Intn(4294967296))
	}
	return id
}

// INTERNAL
// Finds all document(s) that match selection criteria
func (d *Database) findDocuments(selection map[string]interface{}) []*document {
	ind := d.findIndex("ObjectId")
	var docs []*document

	for _, element := range ind.Index {
		docs = append(docs, element.Document)
	}

	for s, i := range selection {
		done := true
		for done {
			for i2, doc := range docs {
				done = false
				if !doc.matches(s, i) {
					done = true
					//fmt.Println("Match:", s, i, (*doc)[s], (*doc)[s] == i)
					if i2 > len(docs) {
						break
					}
					if i2 == len(docs) {
						docs = append([]*document{}, docs[0:i2]...)
						continue
					}
					docs = append(append([]*document{}, docs[0:i2]...), append([]*document{}, docs[i2+1:]...)...)
				}
			}
		}
	}

	return docs
}

// PUBLIC
// for use by others to add to the database
// may change to internal to create naming and regularity among public functions
func (d *Database) Add(data interface{}) Result {
	document, res := convertStruct(data)
	if res.IsError() {
		return *newResult("Failed to convert: "+res.Result(), FAILURE)
	}
	if document == nil {
		panic("Could not convert!")
	}

	document["ObjectId"] = d.generateId()

	d.addDocument(document)

	return *newResult("Added to database", SUCCESS)
}

// PUBLIC
// for use by the end user to interact with the database
func (d *Database) Get(field string, value interface{}) interface{} {
	str, ok := value.(string)

	if !ok {
		return d.findDocumentNoIndex(field, value)
	}

	return d.findDocument(field, str)

}

// PUBLIC
// This is the func that should be used by other packages to get documents
// TODO - make this - iterative indexing
func (d *Database) Read(selection map[string]interface{}, input interface{}) (interface{}, *Result) {
	documents := d.findDocuments(selection)

	fmt.Println(selection)

	if len(documents) > 1 {
		return nil, newResult("Selection Return Multiple Documents!", FAILURE)
	}

	kind := reflect.ValueOf(input).Kind()
	if kind != reflect.Struct {
		return nil, newResult("Output is not Struct!", FAILURE)
	}

	val := reflect.ValueOf(input)

	var output interface{} = reflect.New(val.Type()).Interface()

	fmt.Println(reflect.TypeOf(output), reflect.TypeOf(input), val.Type())

	activeDoc := append([]*document{}, documents[0])
	activeOut := append([]*interface{}{}, &output)

	var lastField []int = []int{0}
	for len(activeOut) != 0 {
		didTerminate := true
		// Get value of current place in output structure
		val = reflect.Indirect(reflect.ValueOf(*activeOut[len(activeOut)-1]))
		fmt.Println("Val", *activeOut[len(activeOut)-1], val.Type(), reflect.Indirect(val).Type())
		for i := lastField[len(lastField)-1]; i < val.NumField(); i++ {
			lastField[len(lastField)-1] = i + 1
			field := val.Field(i)
			typeVal := val.Type().Field(i)
			tag := typeVal.Tag
			name := tag.Get("memdat")
			if name == "" {
				name = typeVal.Name
			}
			if field.Kind() == reflect.Struct {
				var e interface{} = field.Interface()
				var d document = (*activeDoc[len(activeDoc)-1])[name].(document)
				activeOut = append(activeOut, &e)
				activeDoc = append(activeDoc, &d)
				lastField = append(lastField, 0)
				didTerminate = false
				break
			} else if field.Kind() == reflect.Slice {
				// TODO - slice parse
				fmt.Println("Slice Parse:", field.Interface())
				var e interface{} = field.Interface()
				var d interface{} = (*activeDoc[len(activeDoc)-1])[name]
				handleSlice(e, d)
			}
			if field.CanSet() {
				field.Set(reflect.ValueOf((*activeDoc[len(activeDoc)-1])[name]))
			} else {
				fmt.Println("Cannot Set", name, lastField)
			}
		}
		if didTerminate {
			activeOut = append([]*interface{}{}, activeOut[:len(activeOut)-1]...)
			activeDoc = append([]*document{}, activeDoc[:len(activeDoc)-1]...)
			lastField = append([]int{}, lastField[:len(lastField)-1]...)

		}
	}

	/* Dunno abt this part
	currentField := ""
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeVal := val.Type().Field(i)
		tag := typeVal.Tag
		name := tag.Get("memdat")
		if name == "" {
			name = typeVal.Name
		}
		currentField = name
		field.Set(reflect.ValueOf((*activeDoc[len(activeDoc)-1])[currentField]))
		if field.Kind() == reflect.Struct {
			var doc document
			doc = field.Interface().(document)
			activeDoc = append(activeDoc, &doc)

		}
	}
	*/

	return output, newResult("Translated Document", SUCCESS)
}

// INTERNAL
// helper for read method for slices
// TODO - make this
func handleSlice(e interface{}, s interface{}) {
	inVal := reflect.ValueOf(s)
	outVal := reflect.ValueOf(e)
	fmt.Println("VALS-OH!", inVal.Type(), outVal.Type())
	// TODO use makeslice to define specific and known
	for i := 0; i < inVal.NumField(); i++ {
		outVal.Index(i).Set(inVal.Index(i))
	}
	fmt.Println("OUT", outVal.Interface())
}

// PUBLIC
// for use by end user to interact with the database
// TODO - make this
func (d *Database) Operate(selection map[string]interface{}, update map[string]interface{}) Result {
	return *newResult("FUNCTION NOT IMPLEMENTED YET", NO_STATUS)
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
