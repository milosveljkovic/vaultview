package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func colorfulPrint(s string, c tcell.Color) string {
	return fmt.Sprintf("[%s::]%s ", c, s)
}
