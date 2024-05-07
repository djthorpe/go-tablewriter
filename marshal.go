package tablewriter

import (
	"encoding/json"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Marshaller interface {
	Marshal() ([]byte, error)
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Convert any value to a byte array. If quote is true, then the value is
// quoted if it is a string.
func marshal(meta *columnmeta, v any) ([]byte, error) {
	// Returns nil if v is nil
	if v == nil {
		return nil, nil
	}
	// Use marshaller if implemented
	if m, ok := v.(Marshaller); ok {
		return m.Marshal()
	}
	switch v := v.(type) {
	case string:
		// By default, strings are not quoted
		return []byte(v), nil
	default:
		return json.Marshal(v)
	}
}
