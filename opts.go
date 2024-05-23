package tablewriter

import (
	"io"

	// Packages
	"github.com/djthorpe/go-tablewriter/pkg/terminal"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type options struct {
	delim      rune   // Delimiter used to separate fields
	header     bool   // Whether to output a header
	null       string // How the nil value is represented in the output
	timeLayout string // How time values are formatted in the output
	timeLocal  bool   // Whether time values should be printed in local time
	width      int    // Suggested width of the table, including delimiters
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

// Set the suggested width of the table, including delimiters from the terminal width
// if the terminal width is not available, the width is ignored
func OptTerminalWidth(w io.Writer) TableOpt {
	return func(o *options) error {
		if terminal.IsTerminal(w) {
			o.width = terminal.Width(w)
		}
		return nil
	}
}

// Set the suggested width of the table, including delimiters
func OptTableWidth(v int) TableOpt {
	return func(o *options) error {
		if v > 0 {
			o.width = v
		} else {
			return ErrBadParameter.With("OptTableWidth")
		}
		return nil
	}
}

// Output as CSV
func OptOutputCSV() TableOpt {
	return func(o *options) error {
		o.format = formatCSV
		o.delim = ','
		return nil
	}
}

// Output as TSV
func OptOutputTSV() TableOpt {
	return func(o *options) error {
		o.format = formatCSV
		o.delim = '\t'
		return nil
	}
}

// Output as Text
func OptOutputText() TableOpt {
	return func(o *options) error {
		o.format = formatText
		o.delim = '|'
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

// Set the delimiter between fields, if not set with OptOutput...
func OptDelimiter(v rune) TableOpt {
	return func(o *options) error {
		o.delim = v
		return nil
	}
}

// Set how time values are formatted in the output
func OptTimeLayout(v string, local bool) TableOpt {
	return func(o *options) error {
		o.timeLayout = v
		o.timeLocal = local
		return nil
	}
}
