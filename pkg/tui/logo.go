package tui

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var LogoSmall = []string{
	` _    __            ____ _    ___             `,
	`| |  / /___ ___  __/ / /| |  / (_)__ _      __`,
	`| | / / __ '/ / / / / __/ | / / / _ \ | /| / /`,
	`| |/ / /_/ / /_/ / / /_ | |/ / /  __/ |/ |/ / `,
	`|___/\__,_/\__,_/_/\__/ |___/_/\___/|__/|__/  `,
}

type Logo struct {
	*tview.Flex

	logo, status *tview.TextView
	mx           sync.Mutex
}

func NewLogo() *Logo {
	l := Logo{
		Flex:   tview.NewFlex(),
		logo:   logo(),
		status: status(),
	}
	l.SetDirection(tview.FlexRow)
	l.AddItem(l.logo, 5, 1, false)
	l.AddItem(l.status, 2, 1, false)
	l.refreshLogo()
	l.refreshStatus("", tview.Styles.PrimitiveBackgroundColor)

	return &l
}

func (l *Logo) Logo() *tview.TextView {
	return l.logo
}

func (l *Logo) Status() *tview.TextView {
	return l.status
}

func (l *Logo) Reset() {
	l.status.Clear()
	l.status.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
}

func (l *Logo) Err(msg string) {
	l.refreshStatus(msg, tcell.ColorDarkRed)
}

func (l *Logo) Warn(msg string) {
	l.refreshStatus(msg, tcell.ColorDarkOrange)
}

func (l *Logo) Info(msg string) {
	l.refreshStatus(msg, tcell.ColorGray)
}

func (l *Logo) Success(msg string) {
	l.refreshStatus(msg, tcell.ColorGreen)
}

func (l *Logo) refreshStatus(msg string, color tcell.Color) {
	l.status.Clear()
	l.status.SetBackgroundColor(color)
	for i, s := range msg {
		fmt.Fprintf(l.status, "[::b]%s", string(s))
		if i == 45 {
			fmt.Fprintf(l.status, "\n")
		}
		if i == 46*2-3 {
			fmt.Fprintf(l.status, "...")
		}
	}
}

func (l *Logo) refreshLogo() {
	l.logo.Clear()
	for i, s := range LogoSmall {
		fmt.Fprintf(l.logo, "[::b]%s", s)
		if i+1 < len(LogoSmall) {
			fmt.Fprintf(l.logo, "\n")
		}
	}
}

func logo() *tview.TextView {
	v := tview.NewTextView()
	v.SetWordWrap(false)
	v.SetWrap(false)
	v.SetTextAlign(tview.AlignRight)
	v.SetDynamicColors(true)

	return v
}

func status() *tview.TextView {
	v := tview.NewTextView()
	v.SetWordWrap(false)
	v.SetWrap(false)
	v.SetTextAlign(tview.AlignLeft)
	v.SetDynamicColors(true)

	return v
}
