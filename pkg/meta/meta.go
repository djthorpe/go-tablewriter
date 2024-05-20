package meta

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	"github.com/djthorpe/go-tablewriter/pkg/text"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type meta struct {
	fields []*fieldmeta // The fields of the struct
	typ    reflect.Type // The underlying type
	values []any        // The values of the struct
}

type fieldmeta struct {
	field  reflect.StructField
	key    string   // the underlying field name
	name   string   // the output field name
	index  []int    // the index of the field
	tuples []string // tuples from the tags
}

// Struct metadata interface
type Struct interface {
	// Return the underlying type
	Type() reflect.Type

	// Return field metadata
	Fields() []StructField

	// Return the field values in the correct order
	Values(v any) ([]any, error)
}

type StructField interface {
	// Return the field name
	Name() string

	// Whether the field has a tag ie, Is("omitempty")
	Is(name string) bool

	// Return the tag value for a field
	Tag(name string) string

	// Return tuple value
	Tuple(name string) string

	// Return text format
	TextFormat() text.Format
}

///////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

// Create a new metadata object from a struct value and optional
// set of tags
func New(v any, tags ...string) (Struct, error) {
	meta := new(meta)

	// Set parameters
	if rt, _, err := typeOf(v); err != nil {
		return nil, err
	} else {
		meta.typ = rt
	}

	// Set colummns, values, strings
	meta.fields = asColumns(meta.typ, tags)
	meta.values = make([]any, len(meta.fields))

	// Return success
	return meta, nil
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (meta meta) String() string {
	str := "<meta"
	str += " type=" + fmt.Sprint(meta.Type())
	str += " fields=" + fmt.Sprint(meta.fields)
	return str + ">"
}

func (meta fieldmeta) String() string {
	str := "<field"
	if meta.key != meta.name {
		str += fmt.Sprintf(" key=%q", meta.key)
	}
	if meta.name != "" {
		str += fmt.Sprintf(" name=%q", meta.name)
	}
	str += fmt.Sprintf(" index=%v", meta.index)
	if len(meta.tuples) > 0 {
		str += fmt.Sprintf(" tuples=%q", meta.tuples)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return underlying struct type for the data
func (meta *meta) Type() reflect.Type {
	return meta.typ
}

// Return the number of fields
func (meta *meta) NumField() int {
	return len(meta.fields)
}

// Return the fields
func (meta *meta) Fields() []StructField {
	result := make([]StructField, len(meta.fields))
	for i, f := range meta.fields {
		result[i] = f
	}
	return result
}

// Return the field values in the correct order. The input value
// should be a struct
func (meta *meta) Values(v any) ([]any, error) {
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
	for i, col := range meta.fields {
		fv := rv.FieldByIndex(col.index)
		if !fv.IsValid() {
			return nil, ErrBadParameter.Withf("invalid field %q", col.key)
		}
		meta.values[i] = fv.Interface()
	}

	// Return success
	return meta.values, nil
}

// Return a tag value for a field
func (meta *fieldmeta) Name() string {
	if meta.name != "" {
		return meta.name
	}
	return meta.key
}

// Return a tag value for a field
func (meta *fieldmeta) Tag(name string) string {
	return meta.field.Tag.Get(name)
}

// Return named tuple value as a string, or empty string
func (meta *fieldmeta) Tuple(name string) string {
	name = strings.ToLower(name)
	for _, tuple := range meta.tuples {
		parts := strings.SplitN(tuple, ":", 2)
		if parts[0] == name && len(parts) == 2 {
			return parts[1]
		}
	}
	return ""
}

// Return true if the field has the named tuple
func (meta *fieldmeta) Is(name string) bool {
	name = strings.ToLower(name)
	for _, tuple := range meta.tuples {
		parts := strings.SplitN(tuple, ":", 2)
		if parts[0] == name {
			return true
		}
	}
	return false
}

// Return text format
func (meta *fieldmeta) TextFormat() text.Format {
	a := text.Alignment(0)
	if meta.Is("alignleft") {
		a = text.Left
	} else if meta.Is("alignright") {
		a = text.Right
	}
	w, _ := strconv.ParseInt(meta.Tuple("width"), 10, 16)
	return text.Format{
		Width: int(w),
		Align: a,
		Wrap:  meta.Is("wrap"),
	}
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

// asColumns returns a slice of field metadata for a struct type
func asColumns(rt reflect.Type, tag []string) []*fieldmeta {
	cols := make([]*fieldmeta, 0, rt.NumField())

FOR_LOOP:
	for _, f := range reflect.VisibleFields(rt) {
		// Ignore if anonymous or non-exported field
		if f.Anonymous || !f.IsExported() {
			continue
		}

		// Set column metadata
		meta := &fieldmeta{
			field: f,
			key:   f.Name,
			index: f.Index,
		}

		// Process tags
		for _, tag := range tag {
			if value := f.Tag.Get(tag); tag == "" {
				// No tag
				continue
			} else if value == "-" {
				// Ignore field completely
				continue FOR_LOOP
			} else {
				// Set name if first tuple is not empty
				tuples := strings.Split(value, ",")
				if tuples[0] != "" && meta.name == "" {
					meta.name = tuples[0]
				}
				// Add tuples to list of tuples
				meta.tuples = append(meta.tuples, tuples[1:]...)
			}
		}

		// Append column
		cols = append(cols, meta)
	}
	return cols
}
