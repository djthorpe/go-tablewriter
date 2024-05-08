package text

///////////////////////////////////////////////////////////////////////////////
// TYPES

type opts struct {
	delim      rune
	fieldWidth int // Default width of each field
}

// Opt is a function which can be used to set options on the text output
type Opt func(*opts) error

///////////////////////////////////////////////////////////////////////////////
// OPTIONS

// Set the default width of each field
func OptFieldWidth(width int) Opt {
	return func(o *opts) error {
		o.fieldWidth = width
		return nil
	}
}

// Set the field deliminiter, default is '|'
func OptFieldDelim(delim rune) Opt {
	return func(o *opts) error {
		o.delim = delim
		return nil
	}
}
