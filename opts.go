package tablewriter

///////////////////////////////////////////////////////////////////////////////
// TYPES

type opts struct {
	tag   string // Tag used to get additional struct metadata
	delim rune   // Delimiter used to separate fields
}

// TableOpt is a function which can be used to set options on a table
type TableOpt func(*opts) error

///////////////////////////////////////////////////////////////////////////////
// OPTIONS

// Set the struct field tag which is used to set table options, default is "json"
func OptTag(tag string) TableOpt {
	return func(o *opts) error {
		o.tag = tag
		return nil
	}
}

// Set the field delimiter, default is ',' for CSV and '|' for Text
func OptDelim(delim rune) TableOpt {
	return func(o *opts) error {
		o.delim = delim
		return nil
	}
}
