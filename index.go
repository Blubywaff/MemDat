package main

import (
	"math"
)

type index struct {
	Field string
	Index []indexElement
}

type indexElement struct {
	Document *document
	Value    string
}

// INTERNAL
// Finds index of value in the index
func (i *index) findPlace(value string) int {
	bot := 0
	top := len(i.Index)
	cur := -1
	for bot != top {
		cur = (bot+top)/2 + bot
		if i.Index[cur].Value == value {
			break
		} else if i.Index[cur].Value > value {
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
