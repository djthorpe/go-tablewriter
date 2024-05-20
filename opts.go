package tablewriter

///////////////////////////////////////////////////////////////////////////////
// TYPES

type options struct {
	delim      rune   // Delimiter used to separate fields
	header     bool   // Whether to output a header
	null       string // How the nil value is represented in the output
	timeLayout string // How time values are formatted in the output
	timeLocal  bool   // Whether time values should be printed in local time
	format
}

// The output type
type format uint

// TableOpt is a function which can be used to set options on a table
type TableOpt func(*options) error

///////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	_          format = iota // Default output format
	formatCSV                // Output as CSV
	formatText               // Output as text
)

///////////////////////////////////////////////////////////////////////////////
// OPTIONS

// Set the struct field tag which is used to set table options, default is "json"
func OptHeader() TableOpt {
	return func(o *options) error {
		o.header = true
		return nil
	}
}

// Set the field delimiter, default is ',' for CSV and '|' for Text
func OptFieldDelim(delim rune) TableOpt {
	return func(o *options) error {
		o.delim = delim
		return nil
	}
}

// Output as CSV
func OptOutputCSV() TableOpt {
	return func(o *options) error {
		o.format = formatCSV
		return nil
	}
}

// Output as Text
func OptOutputText() TableOpt {
	return func(o *options) error {
		o.format = formatText
		return nil
	}
}

// Set how the nil value is represented in the output, defaults to "<nil>"
func OptNull(v string) TableOpt {
	return func(o *options) error {
		o.null = v
		return nil
	}
}

// Set how time values are formatted in the output
func OptTimeLayout(v string) TableOpt {
	return func(o *options) error {
		o.timeLayout = v
		return nil
	}
}
