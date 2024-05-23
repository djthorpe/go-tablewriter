package text

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Writer struct {
	opts
	w   io.Writer
	row [][]string
}

// Text Alignment
type Alignment int

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultWidth = 20
	defaultDelim = '|'
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWriter(w io.Writer, opts ...Opt) (*Writer, error) {
	writer := new(Writer)
	writer.w = w

	// Set defaults
	writer.delim = defaultDelim
	if err := OptFormat(Format{Width: defaultWidth, Align: Left, Wrap: false})(&writer.opts); err != nil {
		return nil, err
	}

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
	// Set capacity of row
	if cap(w.row) < len(v) {
		w.row = make([][]string, len(v))
	}

	// Format each value
	maxHeight := 0
	for i, value := range v {
		w.row[i] = format(value, w.fieldFormat(i))
		maxHeight = max(len(w.row[i]), maxHeight)
	}

	// Print out each row
	for y := 0; y < maxHeight; y++ {
		for x := 0; x < len(v); x++ {
			if x == 0 {
				w.w.Write([]byte(string(w.delim)))
			}
			if y < len(w.row[x]) {
				w.w.Write([]byte(w.row[x][y]))
			} else {
				w.w.Write([]byte(format("", w.fieldFormat(x))[0]))
			}
			w.w.Write([]byte(string(w.delim)))
		}
		w.w.Write([]byte("\n"))
	}

	// Return success
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// return the format for a row, falling back to the default as needed
func (w *Writer) fieldFormat(i int) Format {
	def := w.format[-1]
	if f, ok := w.format[i]; ok {
		if f.Align == 0 {
			f.Align = def.Align
		}
		if f.Width == 0 {
			f.Width = def.Width
		}
		if !f.Wrap {
			f.Wrap = def.Wrap
		}
		return f
	}
	return def
}

// format a text value to a given format and return the lines
func format(v string, f Format) []string {
	// Trim spaces from the text and reformat
	v = strings.TrimSpace(v)
	if v != "" {
		v = string(quote(v))
	}

	// Wrap text to width (which needs to be greater than 0)
	if f.Wrap {
		v = runewidth.Wrap(v, f.Width)
	}

	// Split lines, truncate and then fill
	lines := strings.Split(v, "\n")
	for i, line := range lines {
		line = runewidth.Truncate(line, f.Width, "")
		switch f.Align {
		case Left:
			lines[i] = runewidth.FillRight(line, f.Width)
		default:
			lines[i] = runewidth.FillLeft(line, f.Width)
		}
	}
	return lines
}

func quote(v string) string {
	var result string
	for _, r := range v {
		result += escapedRune(r)
	}
	return result
}

func escapedRune(r rune) string {
	if r == utf8.RuneError || !utf8.ValidRune(r) {
		return `\uFFFD`
	}
	if strconv.IsPrint(r) {
		return string(r)
	}
	switch r {
	case '\a':
		return `\a`
	case '\b':
		return `\b`
	case '\f':
		return `\f`
	case '\n':
		return `\n`
	case '\r':
		return `\r`
	case '\t':
		return `\t`
	case '\v':
		return `\v`
	default:
		switch {
		case r < ' ' || r == 0x7f:
			return fmt.Sprintf(`\x%02x`, r)
		default:
			return "TODO"
			/*		case :
						return `\uFFFD`
					case r < 0x10000:
						buf = append(buf, `\u`...)
						for s := 12; s >= 0; s -= 4 {
							buf = append(buf, lowerhex[r>>uint(s)&0xF])
						}
					default:
						buf = append(buf, `\U`...)
						for s := 28; s >= 0; s -= 4 {
							buf = append(buf, lowerhex[r>>uint(s)&0xF])
						}*/
		}
	}
}
