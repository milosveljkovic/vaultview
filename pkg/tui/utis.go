package tui

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
	"vaultview/pkg/constants"

	"github.com/gdamore/tcell/v2"
)

const dateFormat = "Jan 2, 2006 3:04 PM"

func colorfulPrint(s string, c tcell.Color) string {
	return fmt.Sprintf("[%s::]%s[::-]", c, s)
}

func formatDate(s string) string {
	layout := time.RFC3339Nano
	t, err := time.Parse(layout, s)
	if err != nil {
		return constants.NAValue
	}
	return t.Format(dateFormat)
}

func getHash(s string) string {
	hash := md5.Sum([]byte(s))
	checksum := hex.EncodeToString(hash[:])
	return checksum
}
