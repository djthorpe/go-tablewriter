package tablewriter

import (
	"fmt"
	"reflect"
	"strings"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type tablemeta struct {
	opts
	cols []*columnmeta // The columns for the table
	typ  reflect.Type  // The underlying type
}

type columnmeta struct {
	Key    string   // the underlying field name
	Name   string   // the output field name
	Index  []int    // the index of the field
	Order  int      // the order of the field, or -1 if the column is to be supressed
	Tuples []string // the tuples from the tag
}

type TableMeta interface {
	// Return the underlying type
	Type() reflect.Type

	// Return the number of output columns
	NumField() int
}

///////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

func (writer *TableWriter) NewMeta(v any, opts ...TableOpt) (*tablemeta, error) {
	meta := new(tablemeta)

	// Set parameters
	if rt, _, err := typeOf(v); err != nil {
		return nil, err
	} else {
		meta.typ = rt
	}

	// Set default options
	meta.opts.tag = defaultTagName

	// Set global options
	for _, opt := range writer.opts {
		if err := opt(&meta.opts); err != nil {
			return nil, err
		}
	}
	// Set local options
	for _, opt := range opts {
		if err := opt(&meta.opts); err != nil {
			return nil, err
		}
	}

	// Set colummns
	meta.cols = asColumns(meta.typ, meta.opts.tag)

	// Return success
	return meta, nil
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (meta tablemeta) String() string {
	str := "<table"
	str += " type=" + fmt.Sprint(meta.Type())
	str += " columns=" + fmt.Sprint(meta.cols)
	return str + ">"
}

func (meta columnmeta) String() string {
	str := "<column"
	if meta.Key != meta.Name {
		str += fmt.Sprintf(" key=%q", meta.Key)
	}
	str += fmt.Sprintf(" name=%q", meta.Name)
	str += fmt.Sprintf(" index=%v", meta.Index)
	str += fmt.Sprintf(" order=%v", meta.Order)
	if len(meta.Tuples) > 0 {
		str += fmt.Sprintf(" tuples=%q", meta.Tuples)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (meta *tablemeta) Type() reflect.Type {
	return meta.typ
}

func (meta *tablemeta) NumField() int {
	return len(meta.cols)
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Returns the type of a value, which is either a slice of structs,
// an array of structs or a single struct. Returns an error if the
// type cannot be determined. If the type is a slice or array, then
// the element type is returned, with the second argument as true.
func typeOf(v any) (reflect.Type, bool, error) {
	// Check parameters
	if v == nil {
		return nil, false, ErrBadParameter.With("nil value")
	}
	rt := reflect.TypeOf(v)
	isSlice := false
	if rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
		rt = rt.Elem()
		isSlice = true
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return nil, false, ErrBadParameter.With("NewTableMeta: not a struct")
	}
	// Return success
	return rt, isSlice, nil
}

// asColumns returns a slice of column metadata for a struct type
func asColumns(rt reflect.Type, tag string) []*columnmeta {
	cols := make([]*columnmeta, 0, rt.NumField())
	order := 0
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		// Ignore private fields
		if f.PkgPath != "" {
			continue
		}

		// Set column metadata
		meta := &columnmeta{
			Key:   f.Name,
			Name:  f.Name,
			Index: f.Index,
		}

		// Obtain tag information from "writer" tag
		if tag := f.Tag.Get(tag); tag != "" {
			tuples := strings.Split(tag, ",")

			// Ignore field if tag is "-"
			if tag == "-" {
				continue
			}

			// Set column output order
			meta.Order = order
			order++

			// Set name if first tuple is not empty
			if tuples[0] != "" && meta.Name == "" {
				meta.Name = tuples[0]
			}

			// Add tuples
			meta.Tuples = append(meta.Tuples, tuples[1:]...)
		}

		// Append column
		cols = append(cols, meta)
	}
	return cols
}
