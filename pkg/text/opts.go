package text

///////////////////////////////////////////////////////////////////////////////
// TYPES

type opts struct {
	delim  rune
	format map[int]Format
}

// Format defines the format of a field
type Format struct {
	// The maximum width of the field
	Width int

	// The alignment of the field (Left or Right)
	Align Alignment

	// Whether to wrap text
	Wrap bool
}

// Opt is a function which can be used to set options on the text output
type Opt func(*opts) error

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	Left Alignment = iota
	Right
)

///////////////////////////////////////////////////////////////////////////////
// OPTIONS

// Set the text format for fields. If fields is omitted, then the format is
// set for all fields
func OptFormat(format Format, fields ...int) Opt {
	return func(o *opts) error {
		if o.format == nil {
			o.format = make(map[int]Format)
		}
		if len(fields) == 0 {
			// Set default format
			o.format[-1] = format
		}
		for _, field := range fields {
			if field >= 0 {
				o.format[field] = format
			}
		}
		return nil
	}
}

// Set the field delimiter, default is '|'
func OptDelim(delim rune) Opt {
	return func(o *opts) error {
		o.delim = delim
		return nil
	}
}
