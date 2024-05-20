package tablewriter

import (
	"encoding/json"
	"reflect"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Marshaller interface {
	Marshal() ([]byte, error)
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Convert any value to a byte array
func marshal(v any, timeLayout string, timeLocal bool) ([]byte, error) {
	// Check for nil
	if v == nil || (reflect.TypeOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil()) {
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
	case time.Time:
		// Return nil if zero time, else return formatted time
		if v.IsZero() {
			return nil, nil
		}
		if timeLocal {
			v = v.Local()
		}
		return []byte(v.Format(timeLayout)), nil
	default:
		return json.Marshal(v)
	}
}
