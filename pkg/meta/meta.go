package meta

import (
	"fmt"
	"reflect"
	"strings"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type meta struct {
	typ    reflect.Type // The underlying type
	fields []*fieldmeta // The fields of the struct
	values []any        // The values of the struct
}

type fieldmeta struct {
	field  reflect.StructField
	key    string   // the underlying field name
	name   string   // the output field name
	index  []int    // the index of the field
	tuples []string // tuples from the tags
	omit   bool     // field output should be omitted
}

// The struct metadata
type Struct interface {
	// Return the underlying type
	Type() reflect.Type

	// Return field metadata
	Fields() []Field

	// Return the field values in the correct order
	Values(v any) ([]any, error)
}

// The field metadata
type Field interface {
	// Return the field name
	Name() string

	// Return the underlying type (dereferencing pointers)
	Type() reflect.Type

	// Return the index of the field
	Index() []int

	// Whether the field has a tag ie, Is("omitempty")
	Is(name string) bool

	// Return the tag value for a field
	Tag(name string) string

	// Return tuple value
	Tuple(name string) string

	// Return omit flag
	Omit() bool

	// Set omit flag
	SetOmit(bool)
}

///////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

// Create a new metadata object from a struct value and optional
// set of tags
func New(v any, tags ...string) (Struct, error) {
	if rt, _, err := typeOf(v); err != nil {
		return nil, err
	} else {
		return NewType(rt, tags...)
	}
}

// Create a new metadata object from a reflect.Type and optional
// set of tags
func NewType(rt reflect.Type, tags ...string) (Struct, error) {
	meta := new(meta)

	// Check type
	if rt.Kind() != reflect.Struct {
		return nil, ErrBadParameter.With("NewType: not a struct")
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
	str += fmt.Sprintf(" kind=%q", meta.field.Type.Kind())
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return underlying struct type for the data
func (meta *meta) Type() reflect.Type {
	return meta.typ
}

// Return the number of fields which are not omitted
func (meta *meta) NumField() int {
	c := 0
	for _, f := range meta.fields {
		if !f.omit {
			c++
		}
	}
	return c
}

// Return the fields
func (meta *meta) Fields() []Field {
	result := make([]Field, 0, len(meta.fields))
	for _, f := range meta.fields {
		if !f.omit {
			result = append(result, f)
		}
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

	// Create a  slice of values
	i := 0
	for _, f := range meta.fields {
		if f.omit {
			continue
		}
		fv := rv.FieldByIndex(f.index)
		if !fv.IsValid() {
			return nil, ErrBadParameter.Withf("invalid field %q", f.key)
		}
		meta.values[i] = fv.Interface()
		i++
	}

	// Return success
	return meta.values[:i], nil
}

// Return a tag value for a field
func (meta *fieldmeta) Name() string {
	if meta.name != "" {
		return meta.name
	}
	return meta.key
}

// Return the index of the field
func (meta *fieldmeta) Index() []int {
	return meta.index
}

// Return the type of field (dereferencing pointers)
func (meta *fieldmeta) Type() reflect.Type {
	result := meta.field.Type
	if result.Kind() == reflect.Ptr {
		result = result.Elem()
	}
	return result
}

// Return the omit flag
func (meta *fieldmeta) Omit() bool {
	return meta.omit
}

// Set the omit flag
func (meta *fieldmeta) SetOmit(v bool) {
	meta.omit = v
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
		if t, err := sliceTypeOf(v); err != nil {
			return nil, false, err
		} else {
			rt = t
			isSlice = true
		}
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

// Returns the type of a slice or array, given that the slice might be an
// interface to many different value types
func sliceTypeOf(v any) (reflect.Type, error) {
	// If it's not an interface, just return the type
	t := reflect.TypeOf(v).Elem()
	if t.Kind() != reflect.Interface {
		return t, nil
	}

	// If it's an interface make sure every element in the slice or array
	// is the same type
	rv := reflect.ValueOf(v)
	if rv.Len() == 0 {
		return nil, ErrBadParameter.With("NewTableMeta: empty []", rv.Type())
	}
	rt := rv.Index(0).Elem().Type()
	for i := 1; i < rv.Len(); i++ {
		if rt != rv.Index(i).Elem().Type() {
			return nil, ErrBadParameter.With("NewTableMeta: []interface with different types")
		}
	}

	// Return (hopefully) the struct type
	return rt, nil
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

		// Set column metadata, default each column to be omitted
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
