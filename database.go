package main

import (
	"fmt"
	"math/rand"
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

func (i Index) contains(value string) bool {
	return i.Index[i.findPlace(value)].Value == value
}

func (i Index) add(document *Document) bool {
	value := (*document)[i.Field].(string)
	place := i.findPlace(value)
	i.Index = append(append(i.Index[0:place], IndexElement{document, value}), i.Index[place:]...)
	return true
}

func (i Index) findDocument(value string) *Document {
	return i.Index[i.findPlace(value)].Document
}

func (d Database) hasIndex(field string) bool {
	return d.findIndexIndex(field) != -1
}

func (d Database) findDocumentById(objectId string) *Document {
	return d.findIndex("ObjectId").findDocument(objectId)
}

func (d Database) findIndex(field string) *Index {
	ii := d.findIndexIndex(field)
	if ii == -1 {
		return nil
	}
	return &d.Indexes[ii]
}

func (d Database) findIndexIndex(field string) int {
	for i, index := range d.Indexes {
		if index.Field == field {
			return i
		}
	}
	return -1
}

func (d Database) addIndex(field string) bool {
	if d.hasIndex(field) {
		return false
	}
	d.Indexes = append(d.Indexes, Index{field, make([]IndexElement, len(d.Documents))})
	for _, document := range d.Documents {
		d.addDocumentToIndex(field, &document)
	}
	return true
}

func (d Database) addDocumentToIndex(field string, document *Document) bool {
	if !d.hasIndex(field) {
		return false
	}
	return d.findIndex(field).add(document)
}

func (d Database) addDocument(document Document) {
	d.Documents = append(d.Documents, document)
	for s, _ := range document {
		d.addDocumentToIndex(s, &document)
	}
}

func (d Database) add(data interface{}) {
	document := Document{"ObjectId": d.generateId()}
	d.addDocument(document)
}

func (d Database) generateId() string {
	id := ""
	for id == "" || d.findIndex("ObjectId").contains(id) {
		id = fmt.Sprintf("%x", rand.Intn(4294967296))
	}
	return id
}

func newDatabase() *Database {
	database := Database{
		make([]Document, 0),
		make([]Index, 0),
	}
	database.addIndex("ObjectId")
	return &database
}
