package main

import (
	"fmt"
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
		cur = (top-bot)/2 + bot
		if top <= bot {
			break
		} else if i.Index[cur].Value == value {
			break
		} else if i.Index[cur].Value > value {
			top = cur - 1
		} else {
			bot = cur + 1
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
func (i *index) add(document *document) Result {
	value := (*document)[i.Field].(string)
	if i.contains(value) {
		return *newResult("Only Uniques are allowed", FAILURE)
	}
	place := int(math.Max(float64(i.findPlace(value)), 0))
	fmt.Println("pint:", place, i.Index[0:place], i.Index[place:])
	fmt.Println("1:", i.Index)
	temp := append((*i).Index[0:place], indexElement{document, value})
	fmt.Println("a:", temp, i.Index)
	i.Index = append(append(i.Index[0:place], indexElement{document, value}), i.Index[place:]...)
	fmt.Println("2:", i.Index)
	return *newResult("Added document to index", SUCCESS)
}

// INTERNAL
// gets document which has the specific value
// TODO needs improvement for robust against findPlace
func (i *index) findDocument(value string) *document {
	return i.Index[i.findPlace(value)].Document
}

// INTERNAL
// Removes the document reference from the index
func (i *index) removeDocument(dptr *document) Result {
	for j := 0; j < len(i.Index); j++ {
		if i.Index[j].Document == dptr {
			i.Index = append(i.Index[0:j], i.Index[j+1:]...)
			return *newResult("Removed Document: "+(*dptr)["ObjectId"].(string), SUCCESS)
		}
	}
	return *newResult("Did not find Document in index: "+i.Field, FAILURE)
}
