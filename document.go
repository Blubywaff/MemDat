package main

// INTERNAL
// used to store data in the database
// should not be used outside this pkg
type document map[string]interface{}

// PUBLIC
// Outer layer field should always be the operation type
//  with another Operation being the value
type Operation struct {
	Field string
	Value map[string]string
}

// INTERNAL
// Used by database to select documents
func (d *document) matches(field string, value interface{}) bool {
	return (&d)[field] == value
}
