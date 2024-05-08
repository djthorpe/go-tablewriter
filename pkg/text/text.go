package text

import (
	"io"
	"strings"

	"github.com/mattn/go-runewidth"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Writer struct {
	opts
	w io.Writer
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWriter(w io.Writer, opts ...Opt) (*Writer, error) {
	writer := new(Writer)
	writer.w = w

	// Set defaults
	writer.delim = '|'
	writer.fieldWidth = 20

	// Set options
	for _, opt := range opts {
		if err := opt(&writer.opts); err != nil {
			return nil, err
		}
	}

	// Return success
	return writer, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (w *Writer) Write(v []string) error {
	for _, value := range v {
		field := fillFieldRow(value, 0, w.fieldWidth)
		w.w.Write([]byte(string(w.delim)))
		w.w.Write([]byte(field))
	}
	w.w.Write([]byte(string(w.delim)))
	w.w.Write([]byte("\n"))

	// Return success
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// fillFieldRow will choose a line in the string y and crop it to the field width x
// filling remaining space with space rune
func fillFieldRow(value string, row, width int) string {
	s := runewidth.Truncate(value, width, "")
	return runewidth.FillRight(rowInString(s, row), width)
}

func rowInString(value string, row int) string {
	lines := strings.Split(value, "\n")
	if row >= 0 && row < len(lines) {
		return lines[row]
	} else {
		return ""
	}
}
