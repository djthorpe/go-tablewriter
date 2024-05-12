package terminal

import (
	"io"
	"os"

	"golang.org/x/term"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// IsTerminal returns true if the file descriptor is a terminal
func IsTerminal(w io.Writer) bool {
	if fd := fileDescriptor(w); fd < 0 {
		return false
	} else {
		return term.IsTerminal(fd)
	}
}

// Width returns the width of the terminal, or zero
func Width(w io.Writer) int {
	if fd := fileDescriptor(w); fd < 0 {
		return 0
	} else if width, _, err := term.GetSize(fd); err != nil {
		return 0
	} else {
		return width
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func fileDescriptor(w io.Writer) int {
	if fh, ok := w.(*os.File); ok {
		return int(fh.Fd())
	} else {
		return -1
	}
}
