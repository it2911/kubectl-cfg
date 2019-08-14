package printers

import (
	"io"

	"github.com/liggitt/tabwriter"
	. "github.com/juju/ansiterm"
)

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
	tabwriterFlags    = tabwriter.RememberWidths
)

// GetNewTabWriter returns a colorable juju/ansiterm/TabWriter that translates tabbed columns in input into properly aligned text.
func GetNewTabWriter(output io.Writer) *TabWriter {
	return NewTabWriter(output, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
}
