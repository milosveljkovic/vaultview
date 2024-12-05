package tui

import (
	"fmt"
	"strconv"
	"vaultview/pkg/constants"
	"vaultview/pkg/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

type SecretDataView struct {
	*tview.Flex
	tui        *Tui
	list       *List
	secret     *tview.TextView
	currentKey string
	secretName string
	keySecret  map[string]string
}

func NewSecretDataView(tui *Tui) *SecretDataView {
	sdw := &SecretDataView{
		Flex: tview.NewFlex(),
		tui:  tui,
		list: NewList(constants.DefaultTitle, tui),
		// keySecret: make(map[string]string),
	}

	sdw.secret = sdw.initSecret()

	sdw.list.EnableSecText()
	sdw.list.List().SetDoneFunc(func() {
		sdw.list.Clear()
		sdw.tui.TogglePage(constants.ViewSecrets)
	})
	sdw.list.List().SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		sdw.currentKey = mainText
	})
	sdw.AddItem(sdw.list.List(), 0, 3, true)
	sdw.AddItem(sdw.secret, 0, 0, false)
	sdw.defineEvents()
	return sdw
}

func (sdw *SecretDataView) initSecret() *tview.TextView {
	s := tview.NewTextView()
	s.SetTextAlign(tview.AlignLeft)
	s.SetBorder(true)
	s.SetWrap(true)
	s.SetDynamicColors(true)
	s.SetDoneFunc(func(key tcell.Key) {
		sdw.secret.Clear()
		sdw.tui.App.SetFocus(sdw.list.List())
		sdw.ResizeItem(sdw.secret, 0, 0)
		sdw.ResizeItem(sdw.list.List(), 0, 3)
	})
	return s
}

func (sdw *SecretDataView) defineEvents() {
	sdw.list.List().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'x' {
			sdw.revealSecret()
			return nil
		}
		return event
	})
	sdw.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'c' {
			//this can cause clipboard_linux.c:15:10: fatal error: X11/Xlib.h
			//depends on libX11, dnf install libX11-devel
			err := clipboard.Init()
			if err != nil {
				sdw.tui.ShowErrAndContinue(fmt.Errorf("copy error: %s", err.Error()))
			}
			s := sdw.keySecret[sdw.currentKey]
			unquoted, _ := strconv.Unquote(s)
			clipboard.Write(clipboard.FmtText, []byte(unquoted))
			sdw.tui.ShowInfoAndContinue("Copied to clipboard!")
			return nil
		}
		return event
	})
}

func (sdw *SecretDataView) revealSecret() {
	s := sdw.keySecret[sdw.currentKey]
	unquoted, _ := strconv.Unquote(s)
	sdw.secret.SetTitle(fmt.Sprintf(" [Secret: [::b]%v[::-], Key: [::b]%v[::-]] ", sdw.secretName, sdw.currentKey))
	fmt.Fprintf(sdw.secret, "%s", unquoted)
	sdw.ResizeItem(sdw.list.List(), 0, 0)
	sdw.ResizeItem(sdw.secret, 0, 1)
	sdw.tui.App.SetFocus(sdw.secret)
}

func (sdw *SecretDataView) Hydrate(data ...interface{}) error {
	sp := data[0].(string)
	se := data[1].(string)
	secrets, err := sdw.tui.vault.ReadKvSecret(se, sp)
	if err != nil {
		sName := utils.GetChildPath(sp)
		sdw.tui.ShowErrAndContinue(fmt.Errorf("secret '%s' does not exist: %v", sName, err))
		sdw.list.Clear()
		sdw.tui.TogglePageAndRefresh(constants.ViewSecrets)
	}
	sdw.secretName = utils.GetChildPath(sp)
	sdw.list.List().SetTitle(fmt.Sprintf(" [Secret: [::b]%v[::-]] ", sdw.secretName))
	sdw.PopulateList(secrets)
	return nil
}

func (sdw *SecretDataView) PopulateList(secrets map[string]string) {
	sdw.keySecret = make(map[string]string)
	for name, s := range secrets {
		sdw.keySecret[name] = s
		sdw.list.Add(name, constants.Mask, nil)
	}
}
