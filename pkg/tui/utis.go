package tui

import (
	"fmt"
	"time"
	"vaultview/pkg/constants"

	"github.com/gdamore/tcell/v2"
)

const dateFormat = "Jan 2, 2006 3:04 PM"

func colorfulPrint(s string, c tcell.Color) string {
	return fmt.Sprintf("[%s::]%s", c, s)
}

func parseTime(s string) string {
	layout := time.RFC3339Nano // This layout matches the provided format
	t, err := time.Parse(layout, s)
	if err != nil {
		return constants.NAValue
	}
	return t.Format(dateFormat)
}
