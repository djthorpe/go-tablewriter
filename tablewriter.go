package tablewriter

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// A tablewriter object
type TableWriter struct {
	w    io.Writer
	opts []TableOpt
	csv  *csv.Writer
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultTagName = "json"
	defaultNull    = "<nil>"
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
	var result error

	// Create a metadata object which creates an iterator for the data
	meta, err := w.NewMeta(v, opts...)
	if err != nil {
		return err
	}

	// If the format is JSON then create a JSON writer
	if meta.opts.format == FormatCSV {
		w.csv = csv.NewWriter(w.w)
		w.csv.Comma = meta.opts.delim
	}

	// Create an iterator
	iterator, err := NewIterator(v)
	if err != nil {
		return err
	}

	// Write rows
	header := false
	for row := iterator.Next(); row != nil; row = iterator.Next() {
		if !header {
			if meta.header {
				if err := w.writeHeader(meta); err != nil {
					result = errors.Join(result, err)
					break
				}
			}
			header = true
		}
		if err := w.writeRow(meta, row); err != nil {
			result = errors.Join(result, err)
		}
	}

	// Flush CSV
	if w.csv != nil {
		w.csv.Flush()
		if err := w.csv.Error(); err != nil {
			result = errors.Join(result, err)
		}
	}

	// Return any errors
	return result
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (w *TableWriter) writeHeader(meta *tablemeta) error {
	switch {
	case w.csv != nil:
		if err := w.csv.Write(meta.Fields()); err != nil {
			return err
		}
	default:
		return errors.New("unsupported format")
	}

	// Return success
	return nil
}

func (w *TableWriter) writeRow(meta *tablemeta, row any) error {
	switch {
	case w.csv != nil:
		if values, err := meta.StringValues(row); err != nil {
			return err
		} else if err := w.csv.Write(values); err != nil {
			return err
		}
	default:
		return errors.New("unsupported format")
	}

	// Return success
	return nil
}
