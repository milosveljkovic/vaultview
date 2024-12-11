package tui

import (
	"vaultview/pkg/models"

	"github.com/rivo/tview"
)

type HeaderViewI interface {
	View
	Info(msg string)
	Err(msg string)
	Reset()
}

type HeaderView struct {
	*tview.Flex
	tui         *Tui
	logo        *Logo
	infoTable   *Info
	placeholder *tview.TextView
}

func NewHeaderView(tui *Tui) *HeaderView {
	header := &HeaderView{
		Flex:        tview.NewFlex(),
		tui:         tui,
		logo:        NewLogo(),
		placeholder: tview.NewTextView(),
	}
	header.infoTable = NewInfo(tui)

	header.SetDirection(tview.FlexColumn)
	header.AddItem(header.infoTable, 80, 1, false).
		AddItem(header.placeholder, 0, 1, false).
		AddItem(header.logo, 46, 1, false)

	return header
}

func (hw *HeaderView) Hydrate(data ...interface{}) error {
	infoModel, err := models.NewInfo(hw.tui.vault, hw.tui.cfg)
	if err != nil {
		return err
	}
	infoModel.RegisterListener(hw.infoTable)
	infoModel.TriggerInfoChange()
	return nil
}

func (hw *HeaderView) Info(msg string) {
	hw.logo.Info(msg)
}

func (hw *HeaderView) Err(msg string) {
	hw.logo.Err(msg)
}

func (hw *HeaderView) Reset() {
	hw.logo.Reset()
}
