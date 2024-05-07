package tablewriter

import (
	"fmt"
	"io"
	"os"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// A tablewriter object
type TableWriter struct {
	w    io.Writer
	opts []TableOpt
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultTagName = "json"
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// New creates a new table writer object, with options for all subsequent writes
func New(w io.Writer, opts ...TableOpt) *TableWriter {
	self := new(TableWriter)
	self.opts = opts
	if w == nil {
		self.w = os.Stdout
	} else {
		self.w = w
	}

	// Return success
	return self
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Write will output the table to the writer object, applying any options
func (w *TableWriter) Write(v any, opts ...TableOpt) error {
	// Create a metadata object which creates an iterator for the data
	meta, err := w.NewMeta(v, opts...)
	if err != nil {
		return err
	}

	// noop
	fmt.Println(meta)

	// Return success
	return nil
}
