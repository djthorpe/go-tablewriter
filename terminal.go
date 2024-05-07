package tablewriter

// TerminalOpts appends appropriate options for terminal output
// including width of the terminal
/*
func TerminalOpts(w io.Writer) []TableOpt {
	result := []TableOpt{}
	if fh, ok := w.(*os.File); ok {
		if term.IsTerminal(int(fh.Fd())) {
			if width, _, err := term.GetSize(int(fh.Fd())); err == nil {
				if width > 2 {
					result = append(result, OptTextWidth(uint(width)))
				}
			}
		}
	}
	return result
}
*/
