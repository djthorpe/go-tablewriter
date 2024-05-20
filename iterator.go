package tablewriter

import (
	"fmt"
	"reflect"

	"github.com/djthorpe/go-tablewriter/pkg/meta"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type iterator struct {
	slice reflect.Value
	index int
}

// Iterator is an interface for iterating over a slice of struct values
type Iterator interface {
	// Return the number of elements
	Len() int

	// Return the next struct, or nil
	Next() any

	// Reset the iterator to the beginning
	Reset()
}

///////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

// NewIterator returns a new slice iterator object, from a single struct
// value or an array of one or more struct values which are of the same type
func NewIterator(v any) (Iterator, error) {
	self := new(iterator)

	// Get the type
	rt, isSlice, err := meta.TypeOf(v)
	if err != nil {
		return nil, err
	}

	// Set the slice parameter
	if isSlice {
		self.slice = reflect.ValueOf(v)
	} else {
		self.slice = reflect.MakeSlice(reflect.SliceOf(rt), 1, 1)
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			self.slice.Index(0).Set(rv.Elem())
		} else {
			self.slice.Index(0).Set(rv)
		}
	}

	// Set the index parameter
	self.index = 0

	// Return success
	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (i iterator) String() string {
	str := "<iterator"
	str += fmt.Sprint(" len=", i.slice.Len())
	str += fmt.Sprint(" i=", i.index)
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the number of elements
func (i *iterator) Reset() {
	i.index = 0
}

// Return the number of elements
func (i *iterator) Len() int {
	return i.slice.Len()
}

// Return the next struct, or nil
func (i *iterator) Next() any {
	if i.index >= i.slice.Len() {
		return nil
	}
	v := i.slice.Index(i.index).Interface()
	i.index++
	return v
}
