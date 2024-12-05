package tui

import (
	"fmt"
	"net/http"
	"strings"
	"vaultview/pkg/constants"
	"vaultview/pkg/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SecretViewI interface {
	View
	SecretsHardRefresh()
}

type SecretView struct {
	*tview.Flex
	tui                   *Tui
	path                  *tview.TextView
	list                  *List
	cachedSecrets         map[string][]string
	engine, currentSecret string
}

func NewSecretView(tui *Tui) *SecretView {
	sw := &SecretView{
		Flex:          tview.NewFlex(),
		tui:           tui,
		list:          NewList(constants.SecretsTitle, tui),
		path:          path(),
		cachedSecrets: make(map[string][]string),
	}

	sw.SetDirection(tview.FlexRow)
	sw.AddItem(sw.path, 3, 1, false)

	sw.list.List().SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		sw.currentSecret = mainText
	})
	sw.list.List().SetDoneFunc(func() {
		p := sw.getPath()
		sw.list.List().Clear().SetOffset(0, 0)
		if p == "" {
			sw.path.Clear()
			sw.setEngine("")
			sw.tui.TogglePage(constants.ViewSecretEngines)
		} else {
			sw.hydratePreviousSecret(p)
		}
	})
	sw.AddItem(sw.list.List(), 0, 3, true)
	sw.defineEvents()

	return sw
}

// load previous path
func (sw *SecretView) hydratePreviousSecret(p string) {
	parentPath := utils.GetParentPath(p)
	selectChild := utils.GetChildPath(p)
	sw.setPath(parentPath)
	if parentPath != "" {
		sw.list.Hydrate(sw.cachedSecrets[sw.getCachedSecretKey(parentPath)], selectChild)
	} else {
		sw.list.Hydrate(sw.cachedSecrets[""])
	}
}

func (sw *SecretView) defineEvents() {
	sw.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			sw.SelectedSecret()
			return nil
		} else if event.Key() == tcell.KeyCtrlR {
			sw.secretsHardRefresh()
			return nil
		}
		return event
	})
}

func path() *tview.TextView {
	path := tview.NewTextView()
	path.SetTextAlign(tview.AlignLeft)
	path.SetBorder(true)
	return path
}

func (sw *SecretView) getPath() string {
	return strings.TrimPrefix(sw.path.GetText(true), "path: /")
}

func (sw *SecretView) setPath(path string) {
	sw.path.SetText(fmt.Sprintf("path: /%s", path))
}

func (sw *SecretView) setPathTitle(pathTitle string) {
	sw.path.SetTitle(fmt.Sprintf(" [Secret Engine: [::b]%v[::-]] ", pathTitle))
}

func (sw *SecretView) setEngine(engine string) {
	sw.engine = engine
}

func (sw *SecretView) getCachedSecretKey(path string) string {
	return sw.engine + ":" + path
}

func (sw *SecretView) Hydrate(data ...interface{}) error {
	sw.setPath("")
	if engine, ok := data[0].(string); ok {
		sw.setEngine(engine)
		sw.setPathTitle(engine)
	} else {
		return fmt.Errorf("error during type assertion")
	}

	vs, err := sw.tui.vault.ListKvSecrets(sw.engine, "")
	if err != nil {
		return err
	}
	sw.cachedSecrets[""] = vs
	sw.list.Hydrate(sw.cachedSecrets[""])

	return nil
}

func (sw *SecretView) SelectedSecret() {
	p := sw.getPath() + sw.currentSecret
	sePath := sw.getCachedSecretKey(p)
	if strings.HasSuffix(sw.currentSecret, "/") {
		// sw.list.Clear()
		if _, ok := sw.cachedSecrets[sePath]; !ok {
			secrets, err := sw.tui.vault.ListKvSecrets(sw.engine, p)
			if err != nil {
				sw.tui.ShowErrAndContinue(err)
				sw.secretsHardRefresh()
				return
			}
			sw.cachedSecrets[sePath] = secrets
		}
		sw.setPath(p)
		sw.list.Hydrate(sw.cachedSecrets[sePath])
	} else {
		sw.tui.ShowSecretDataView(p, sw.engine)
	}
}

func (sw *SecretView) secretsHardRefresh() {
	p := sw.getPath()
	go func(path string) {
		sw.tui.App.QueueUpdateDraw(func() {
			secrets, err := sw.tui.vault.ListKvSecrets(sw.engine, path)
			currSEPath := sw.getCachedSecretKey(p)
			if err != nil {
				if sw.tui.vault.IsErrorStatus(err, http.StatusNotFound) {
					sw.tui.ShowErrAndContinue(fmt.Errorf("path '%s' does not exist: %v", p, err))
					seParentPath := sw.getCachedSecretKey(utils.GetParentPath(p))
					seChildPath := sw.getCachedSecretKey(utils.GetChildPath(p))
					delete(sw.cachedSecrets, currSEPath)
					// remove childPath from parentPath list
					sw.cachedSecrets[seParentPath] = utils.RemoveFromSlice(sw.cachedSecrets[seParentPath], seChildPath)
					sw.hydratePreviousSecret(p)
					return
				}
			}
			sw.cachedSecrets[currSEPath] = secrets
			sw.list.Hydrate(sw.cachedSecrets[currSEPath])
		})
	}(p)
	sw.list.Clear()
}

func (sw *SecretView) SecretsHardRefresh() {
	sw.secretsHardRefresh()
}
