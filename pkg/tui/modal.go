package tui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ModalInput struct {
	*tview.Form
	tui          *Tui
	DialogHeight int
	frame        *tview.Frame
	primary      string
	secondary    string
	err          string
	done         func(string, string, bool)
}

func NewModalInput(tui *Tui) *ModalInput {
	form := tview.NewForm()

	m := &ModalInput{
		form,
		tui,
		9,
		tview.NewFrame(form),
		"",
		"",
		"",
		nil,
	}

	m.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor).
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetBorderPadding(0, 0, 0, 0)

	m.frame.SetTitle(fmt.Sprintf(" [%s::]%s ", tcell.ColorWhite, "Vault Configuration"))
	m.frame.SetBorders(0, 0, 1, 0, 0, 0).
		SetBorder(true).
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetBorderPadding(1, 1, 1, 1)

	m.SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	m.SetCancelFunc(func() {
		if m.done != nil {
			m.done("", "", false)
		}
	})

	m.AddButton("OK", func() {
		if m.done != nil {
			m.done(m.primary, m.secondary, true)
		}
	})
	m.AddButton("Cancel", func() {
		if m.done != nil {
			m.done(m.primary, m.secondary, false)
		}
	})

	return m
}

func (m *ModalInput) Init() {
	//todo: refactor this
	if os.Getenv("VAULT_ADDR") != "" {
		m.AddInputField(colorfulPrint("Vault Addr:", tcell.ColorLime), os.Getenv("VAULT_ADDR"), 0, nil, func(text string) {
			m.primary = text
		})
	} else {
		m.AddInputField(colorfulPrint("Vault Addr:", tcell.ColorLime), "", 0, nil, func(text string) {
			m.primary = text
		})
	}
	if os.Getenv("VAULT_TOKEN") != "" {
		m.AddPasswordField(colorfulPrint("Vault Token:", tcell.ColorLime), os.Getenv("VAULT_TOKEN"), 0, '*', func(text string) {
			m.secondary = text
		})
	} else {
		m.AddPasswordField(colorfulPrint("Vault Token:", tcell.ColorLime), "", 0, '*', func(text string) {
			m.secondary = text
		})
	}

	m.SetDoneFunc(func(primText, secText string, success bool) {
		if success {
			m.tui.cfg.UpdateVaultAddr(primText)
			m.tui.InitVault(primText, secText)
			err := m.tui.InitMain()
			if err != nil {
				m.tui.ShowErrAndStop(err)
			}
		} else {
			err := m.tui.InitMain()
			if err != nil {
				m.tui.ShowErrAndStop(err)
			}
		}
	})
}

// SetValue sets the current value in the item
func (m *ModalInput) SetValue(text string, secondary string) {
	m.primary = text
	m.secondary = secondary
	m.Clear(false)
	m.AddInputField("", text, 50, nil, func(text string) {
		if len(text) == 0 {
			text = "(empty)"
		}
		m.primary = text
	})
	m.AddInputField("", secondary, 50, nil, func(text string) {
		m.secondary = text
	})
}

func (m *ModalInput) setPrimText(text string) {
	m.primary = text
}

func (m *ModalInput) setSecText(text string) {
	m.secondary = text
}

// SetDoneFunc sets the done func for this input.
// Will be called with the text of the input and a boolean for OK or cancel button.
func (m *ModalInput) SetDoneFunc(handler func(string, string, bool)) *ModalInput {
	m.done = handler
	return m
}

// Draw draws this primitive onto the screen.
func (m *ModalInput) Draw(screen tcell.Screen) {
	// Calculate the width of this modal.
	buttonsWidth := 50
	screenWidth, screenHeight := screen.Size()
	width := screenWidth / 3
	if width < buttonsWidth {
		width = buttonsWidth
	}
	// width is now without the box border.

	// Set the modal's position and size.
	height := m.DialogHeight
	width += 1
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2
	m.SetRect(x, y, width, height)

	// Draw the frame.
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}
