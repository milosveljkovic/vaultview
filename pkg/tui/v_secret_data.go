package tui

import (
	"fmt"
	"vaultview/pkg/constants"
	"vaultview/pkg/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

type SecretMetadata struct {
	version      string
	created_time string
}

type SecretDataView struct {
	*tview.Flex
	tui                      *Tui
	list                     *List
	secret                   *tview.TextView
	editor                   *tview.TextArea
	currentKey, secretName   string
	secretEng, secretPath    string
	keySecret, editKeySecret map[string]string
	metadata                 SecretMetadata
}

func NewSecretDataView(tui *Tui) *SecretDataView {
	sdw := &SecretDataView{
		Flex: tview.NewFlex(),
		tui:  tui,
		list: NewList(constants.DefaultTitle, tui),
	}

	sdw.secret = sdw.initSecret()
	sdw.editor = sdw.initEditor()

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
	sdw.AddItem(sdw.editor, 0, 0, false)
	sdw.defineEvents()
	return sdw
}

func (sdw *SecretDataView) initSecret() *tview.TextView {
	s := tview.NewTextView()
	s.SetTextAlign(tview.AlignLeft)
	s.SetBorder(true)
	s.SetWrap(true)
	s.SetDynamicColors(true)
	s.SetTitle(fmt.Sprint(" [[::b]Preview Mode[::-]] "))
	s.SetDoneFunc(func(key tcell.Key) {
		sdw.list.List().SetTitle(sdw.getFancyTitle())
		sdw.secret.Clear()
		sdw.tui.App.SetFocus(sdw.list.List())
		sdw.ResizeItem(sdw.secret, 0, 0)
		sdw.ResizeItem(sdw.editor, 0, 0)
		sdw.ResizeItem(sdw.list.List(), 0, 3)
	})
	return s
}

func (sdw *SecretDataView) initEditor() *tview.TextArea {
	s := tview.NewTextArea()
	s.SetBorder(true)
	s.SetBorderColor(tcell.ColorLime)
	s.SetBorderAttributes(tcell.AttrBold)
	s.SetTitle(fmt.Sprint(" [[::b]Edit Mode[::-]] "))
	s.SetWrap(true)
	s.SetChangedFunc(func() {
		sdw.editKeySecret[sdw.currentKey] = sdw.editor.GetText()
	})
	s.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			sdw.list.List().SetTitle(sdw.getFancyTitle())
			sdw.secret.Clear()
			sdw.tui.App.SetFocus(sdw.list.List())
			sdw.ResizeItem(sdw.secret, 0, 0)
			sdw.ResizeItem(sdw.editor, 0, 0)
			sdw.ResizeItem(sdw.list.List(), 0, 3)
			return nil
		} else if event.Key() == tcell.KeyCtrlS {
			sdw.SaveSecret()
			return nil
		} else if event.Key() == tcell.KeyTab {
			// tab is used to switch between the keys (not available in the editor)
			sdw.list.NextItem()
			sdw.activateEditor()
			return nil
		}
		return event
	})
	return s
}

func (sdw *SecretDataView) defineEvents() {
	sdw.list.List().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == constants.Reveal {
			sdw.revealSecret()
			return nil
		} else if event.Rune() == constants.Copy {
			sdw.CopyToClipboard()
			return nil
		} else if event.Rune() == constants.Edit {
			sdw.activateEditor()
			return nil
		} else if event.Key() == tcell.KeyCtrlS {
			sdw.SaveSecret()
			return nil
		}
		return event
	})
	sdw.secret.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == constants.Copy {
			sdw.CopyToClipboard()
			return nil
		} else if event.Rune() == constants.Edit {
			sdw.activateEditor()
			return nil
		} else if event.Key() == tcell.KeyTab {
			// tab is used to switch between the keys (not available in the editor)
			sdw.list.NextItem()
			sdw.revealSecret()
			return nil
		}
		return event
	})
}

func (sdw *SecretDataView) activateEditor() {
	s := sdw.editKeySecret[sdw.currentKey]
	sdw.list.List().SetTitle(sdw.getFancyTitleShort())
	sdw.editor.SetText(s, false)
	sdw.ResizeItem(sdw.list.List(), 0, 1)
	sdw.ResizeItem(sdw.secret, 0, 0)
	sdw.ResizeItem(sdw.editor, 0, 3)
	sdw.tui.App.SetFocus(sdw.editor)
}

func (sdw *SecretDataView) revealSecret() {
	sdw.secret.Clear()
	s := sdw.keySecret[sdw.currentKey]
	sdw.list.List().SetTitle(sdw.getFancyTitleShort())
	fmt.Fprintf(sdw.secret, "%s", s)
	sdw.ResizeItem(sdw.list.List(), 0, 1)
	sdw.ResizeItem(sdw.editor, 0, 0)
	sdw.ResizeItem(sdw.secret, 0, 3)
	sdw.tui.App.SetFocus(sdw.secret)
}

func (sdw *SecretDataView) Hydrate(data ...interface{}) error {
	sdw.secretPath = data[0].(string)
	sdw.secretEng = data[1].(string)
	secrets, metadata, err := sdw.tui.vault.ReadKvSecret(sdw.secretEng, sdw.secretPath)
	if err != nil {
		sName := utils.GetChildPath(sdw.secretPath)
		sdw.tui.ShowStatusAndContinue(fmt.Sprintf("secret '%s' does not exist: %v", sName, err), ErrStatus)
		sdw.list.Clear()
		sdw.tui.TogglePageAndRefresh(constants.ViewSecrets)
	}
	sdw.secretName = utils.GetChildPath(sdw.secretPath)
	sdw.metadata = SecretMetadata{
		version:      metadata["version"],
		created_time: formatDate(metadata["created_time"]),
	}
	sdw.list.List().SetTitle(sdw.getFancyTitle())
	sdw.PopulateList(secrets)
	return nil
}

func (sdw *SecretDataView) CopyToClipboard() {
	err := clipboard.Init()
	if err != nil {
		sdw.tui.ShowStatusAndContinue(fmt.Sprintf("Copy to clipboard error: %s", err.Error()), ErrStatus)
	}
	s := sdw.keySecret[sdw.currentKey]
	clipboard.Write(clipboard.FmtText, []byte(s))
	sdw.tui.ShowStatusAndContinue("Copied to clipboard", InfoStatus)
}

func (sdw *SecretDataView) SaveSecret() {
	hasChanged := false
	for k, v := range sdw.keySecret {
		currentHash := getHash(v)
		newHash := getHash(sdw.editKeySecret[k])
		if currentHash != newHash {
			hasChanged = true
			break
		}
	}
	if hasChanged {
		sdwKeySecretAny := make(map[string]any)
		for k, v := range sdw.editKeySecret {
			sdwKeySecretAny[k] = v
		}
		if err := sdw.tui.vault.WriteKv2Secret(sdw.secretEng, sdw.secretPath, sdwKeySecretAny); err != nil {
			sdw.tui.ShowStatusAndContinue(err.Error(), ErrStatus)
		} else {
			sdw.tui.ShowStatusAndContinue(fmt.Sprintf("Secret '%s' updated successfuly", sdw.secretName), SuccessStatus)
			sdw.tui.App.SetFocus(sdw.list.List())
			sdw.secret.Clear()
			sdw.list.Clear()
			sdw.Hydrate(sdw.secretPath, sdw.secretEng)
			sdw.ResizeItem(sdw.secret, 0, 0)
			sdw.ResizeItem(sdw.editor, 0, 0)
			sdw.ResizeItem(sdw.list.List(), 0, 3)
		}
	} else {
		sdw.tui.ShowStatusAndContinue("Nothing to save...", InfoStatus)
	}
}

func (sdw *SecretDataView) getFancyTitle() string {
	return fmt.Sprintf(" [%v[::b] %v[::-], %s[::b] %v[::-], %s[::b] %v[::-]] ", "Secret:", sdw.secretName, "Ver:", sdw.metadata.version, "Created:", sdw.metadata.created_time)
}

func (sdw *SecretDataView) getFancyTitleShort() string {
	return fmt.Sprintf(" [%v[::b] %v[::-], %s[::b] %v[::-]] ", "Secret:", sdw.secretName, "Ver:", sdw.metadata.version)
}

func (sdw *SecretDataView) PopulateList(secrets map[string]string) {
	sdw.keySecret = make(map[string]string)
	sdw.editKeySecret = make(map[string]string)
	for name, s := range secrets {
		sdw.keySecret[name] = s
		sdw.editKeySecret[name] = s
		sdw.list.Add(name, constants.Mask, nil)
	}
}
