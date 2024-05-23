package tablewriter

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"time"

	// Packages
	meta "github.com/djthorpe/go-tablewriter/pkg/meta"
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
	row  []string
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultNull       = "<nil>"
	defaultTimeLayout = time.RFC1123
	defaultTimeLocal  = false
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

	// Create an iterator
	iterator, err := meta.NewIterator(v)
	if err != nil {
		return err
	}

	// Create a metadata object which creates an iterator for the data
	meta, err := meta.New(v, "writer", "json")
	if err != nil {
		return err
	}

	// Options processing
	var o options
	o.format = formatCSV
	o.delim = ','
	o.null = defaultNull
	o.timeLayout = defaultTimeLayout
	o.timeLocal = defaultTimeLocal
	for _, opt := range append(w.opts, opts...) {
		if err := opt(&o); err != nil {
			return err
		}
	}

	// Check for zeroed-data columns - initalize the "notomit"
	// slice to false, and then iterate over the rows to see if
	// any columns are not zeroed, flagging them as "notomit"
	fields := meta.Fields()
	notomit := make([]bool, len(fields))
	for row := iterator.Next(); row != nil; row = iterator.Next() {
		values, err := meta.Values(row)
		if err != nil {
			return err
		}
		for i, value := range values {
			if notomit[i] {
				continue
			}
			if value == nil {
				continue
			}
			if reflect.ValueOf(value).IsZero() {
				continue
			}
			notomit[i] = true
		}
	}
	iterator.Reset()

	// Set omit flags based on the notomit slice and the "omitempty" tag
	for i, field := range fields {
		if field.Is("omitempty") && !notomit[i] {
			field.SetOmit(true)
		} else {
			field.SetOmit(false)
		}
	}

	// Create the writer object based on the format required
	switch o.format {
	case formatCSV:
		w.csv = csv.NewWriter(w.w)
		w.csv.Comma = o.delim
	case formatText:
		opts := []text.Opt{
			text.OptDelim(o.delim),
		}
		for i, field := range meta.Fields() {
			if textFormat := textFormat(field); textFormat.Width > 0 || textFormat.Align != 0 || textFormat.Wrap {
				opts = append(opts, text.OptFormat(textFormat, i))
			}
		}
		if o.width > 0 {
			fmt.Println("TODO: Set width", o.width)
		}
		if writer, err := text.NewWriter(w.w, opts...); err != nil {
			return err
		} else {
			w.text = writer
		}
	default:
		return errUnsupportedFormat
	}

	// Write rows
	header := false
	for row := iterator.Next(); row != nil; row = iterator.Next() {
		if !header {
			if o.header {
				if err := w.writeHeader(o.format, meta); err != nil {
					result = errors.Join(result, err)
					break
				}
			}
			header = true
		}
		if err := w.writeRow(&o, meta, row); err != nil {
			result = errors.Join(result, err)
		}
	}

	// Flush
	switch o.format {
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

func textFormat(field meta.Field) text.Format {
	var result text.Format

	// Wrap
	if field.Is("wrap") {
		result.Wrap = true
	}

	// Alignment
	switch {
	case field.Is("left"):
		result.Align = text.Left
	case field.Is("right"):
		result.Align = text.Right
	}

	// Width
	if field.Is("width") {
		if w, err := strconv.ParseInt(field.Tuple("width"), 10, 16); err == nil {
			result.Width = int(w)
		}
	}
	return result
}

func (w *Writer) writeHeader(f format, meta meta.Struct) error {
	fields := meta.Fields()
	w.row = make([]string, len(fields))
	for i, field := range fields {
		w.row[i] = field.Name()
	}

	// Write header row
	switch f {
	case formatCSV:
		if err := w.csv.Write(w.row); err != nil {
			return err
		}
	case formatText:
		if err := w.text.Write(w.row); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

func (w *Writer) writeRow(o *options, meta meta.Struct, row any) error {
	values, err := meta.Values(row)
	if err != nil {
		return err
	}

	// Convert values to []string
	if len(w.row) != len(values) {
		w.row = make([]string, len(values))
	}

	// Marshal values
	var result error
	for i, v := range values {
		if cell, err := marshal(v, o.timeLayout, o.timeLocal); err != nil {
			result = errors.Join(result, err)
		} else if cell == nil {
			w.row[i] = o.null
		} else {
			w.row[i] = string(cell)
		}
	}
	if result != nil {
		return result
	}

	// Write row
	switch o.format {
	case formatCSV:
		if err := w.csv.Write(w.row); err != nil {
			return err
		}
	case formatText:
		if err := w.text.Write(w.row); err != nil {
			return err
		}
	}

	// Return success
	return nil
}
