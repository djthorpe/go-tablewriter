package tablewriter

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"time"

	// Packages
	text "github.com/djthorpe/go-tablewriter/pkg/text"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// A writer object which can write table data to an io.Writer
type Writer struct {
	w    io.Writer
	opts []TableOpt
	csv  *csv.Writer
	text *text.Writer
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultTagName    = "json"
	defaultNull       = "<nil>"
	defaultTimeLayout = time.RFC1123
)

var (
	errUnsupportedFormat = errors.New("unsupported output format")
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// New creates a new Writer object, with options for all subsequent writes
func New(w io.Writer, opts ...TableOpt) *Writer {
	self := new(Writer)
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

// Output will return the underlying io.Writer object
func (w *Writer) Output() io.Writer {
	return w.w
}

// Write the table to output, applying any options which override to the
// options passed to the New method
func (w *Writer) Write(v any, opts ...TableOpt) error {
	var result error

	// Create a metadata object which creates an iterator for the data
	meta, err := w.NewMeta(v, opts...)
	if err != nil {
		return err
	}

	// If the format is JSON then create a JSON writer
	switch meta.(*tablemeta).opts.format {
	case formatCSV:
		w.csv = csv.NewWriter(w.w)
		w.csv.Comma = meta.(*tablemeta).opts.delim
	case formatText:
		if writer, err := text.NewWriter(w.w); err != nil {
			return err
		} else {
			w.text = writer
		}
	default:
		return errUnsupportedFormat
	}

	// Create an iterator
	iterator, err := NewIterator(v)
	if err != nil {
		return err
	}

	// TODO: Check for zeroed-data columns
	//for row := iterator.Next(); row != nil; row = iterator.Next() {
	//	if err := meta.CheckZero(row); err != nil {
	//		result = errors.Join(result, err)
	//	}
	//}
	//iterator.Reset()

	// Write rows
	header := false
	for row := iterator.Next(); row != nil; row = iterator.Next() {
		if !header {
			if meta.(*tablemeta).header {
				if err := w.writeHeader(meta.(*tablemeta)); err != nil {
					result = errors.Join(result, err)
					break
				}
			}
			header = true
		}
		if err := w.writeRow(meta.(*tablemeta), row); err != nil {
			result = errors.Join(result, err)
		}
	}

	// Flush
	switch meta.(*tablemeta).opts.format {
	case formatCSV:
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

func (w *Writer) writeHeader(meta *tablemeta) error {
	switch meta.opts.format {
	case formatCSV:
		if err := w.csv.Write(meta.Fields()); err != nil {
			return err
		}
	case formatText:
		if err := w.text.Write(meta.Fields()); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

func (w *Writer) writeRow(meta *tablemeta, row any) error {
	switch meta.opts.format {
	case formatCSV:
		if values, err := meta.StringValues(row); err != nil {
			return err
		} else if err := w.csv.Write(values); err != nil {
			return err
		}
	case formatText:
		if values, err := meta.StringValues(row); err != nil {
			return err
		} else if err := w.text.Write(values); err != nil {
			return err
		}
	}

	// Return success
	return nil
}
