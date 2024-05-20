package tablewriter

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	// Package imports
	text "github.com/djthorpe/go-tablewriter/pkg/text"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type tablemeta struct {
	opts
	cols    []*columnmeta // The columns for the table
	typ     reflect.Type  // The underlying type
	values  []any         // The current row
	strings []string      // The current row as strings
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

	// Return the field names
	Fields() []string
}

///////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

// Create a new table metadata object from a struct value and optional
// formatting options
func (writer *Writer) NewMeta(v any, opts ...TableOpt) (TableMeta, error) {
	meta := new(tablemeta)

	// Set parameters
	if rt, _, err := typeOf(v); err != nil {
		return nil, err
	} else {
		meta.typ = rt
	}

	// Set default options
	meta.opts.tag = defaultTagName
	meta.opts.null = defaultNull
	meta.opts.format = formatCSV
	meta.opts.delim = ','
	meta.opts.timeLayout = defaultTimeLayout

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

	// Set colummns, values, strings
	meta.cols = asColumns(meta.typ, meta.opts.tag)
	meta.values = make([]any, len(meta.cols))
	meta.strings = make([]string, len(meta.cols))

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

// Return underlying struct type for the data
func (meta *tablemeta) Type() reflect.Type {
	return meta.typ
}

// Return the number of fields or columns
func (meta *tablemeta) NumField() int {
	return len(meta.cols)
}

// Return the field names
func (meta *tablemeta) Fields() []string {
	for i, col := range meta.cols {
		meta.strings[i] = col.Name
	}
	return meta.strings
}

// Return the field formats
func (meta *tablemeta) Formats() []text.Format {
	result := make([]text.Format, len(meta.cols))
	for i, col := range meta.cols {
		result[i] = text.Format{
			Wrap:  col.Bool("wrap"),
			Width: col.Int("width"),
		}
		if col.Bool("right") {
			result[i].Align = text.Right
		} else {
			result[i].Align = text.Left
		}
	}
	return result
}

// Return the field values in the correct order. The input value should be
// a struct
func (meta *tablemeta) Values(v any) ([]any, error) {
	if v == nil {
		return nil, ErrBadParameter.With("nil value")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Type() != meta.typ {
		return nil, ErrBadParameter.Withf("expected type %q", meta.typ)
	}

	// Create a slice of values
	for i, col := range meta.cols {
		fv := rv.FieldByIndex(col.Index)
		if !fv.IsValid() {
			return nil, ErrBadParameter.Withf("invalid field %q", col.Key)
		}
		meta.values[i] = fv.Interface()
	}

	// Return success
	return meta.values, nil
}

// Return the field values as marshalled strings in the correct order.
// The input value should be a struct
func (meta *tablemeta) StringValues(v any) ([]string, error) {
	var result error
	values, err := meta.Values(v)
	if err != nil {
		return nil, err
	}
	for i, value := range values {
		if bytes, err := marshal(meta.cols[i], value, meta.timeLayout); err != nil {
			result = errors.Join(result, err)
		} else if bytes == nil {
			meta.strings[i] = meta.opts.null
		} else {
			meta.strings[i] = string(bytes)
		}
	}

	// Return any errors
	return meta.strings, result
}

// Return true if the column has the named tuple
func (meta *columnmeta) Bool(name string) bool {
	name = strings.ToLower(name)
	for _, tuple := range meta.Tuples {
		parts := strings.SplitN(tuple, ":", 2)
		if parts[0] == name {
			return true
		}
	}
	return false
}

// Return named tuple value as an int, or zero
func (meta *columnmeta) Int(name string) int {
	name = strings.ToLower(name)
	for _, tuple := range meta.Tuples {
		parts := strings.SplitN(tuple, ":", 2)
		if parts[0] == name && len(parts) == 2 {
			if v, err := strconv.ParseInt(parts[1], 10, 0); err == nil {
				return int(v)
			}
		}
	}
	return 0
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
	for _, f := range reflect.VisibleFields(rt) {
		// Ignore if anonymous field
		if f.Anonymous {
			continue
		}
		// Ignore unexported fields
		if !f.IsExported() {
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
			// Ignore field if tag is "-"
			if tag == "-" {
				continue
			}

			// Set column output order
			meta.Order = order
			order++

			// Set name if first tuple is not empty
			tuples := strings.Split(tag, ",")
			if tuples[0] != "" {
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
