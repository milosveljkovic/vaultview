package tui

import (
	"vaultview/pkg/constants"

	"github.com/rivo/tview"
)

type SecretEngineView struct {
	*tview.Flex
	tui  *Tui
	list *List
}

func NewSecretEngineView(tui *Tui) *SecretEngineView {
	secretView := &SecretEngineView{
		Flex: tview.NewFlex(),
		tui:  tui,
		list: NewList(constants.SecretEnginesTitle, tui),
	}

	secretView.AddItem(secretView.list.List(), 0, 3, true)

	return secretView
}

func (sew *SecretEngineView) Hydrate(data ...interface{}) error {
	se, err := sew.tui.vault.ReadSecretEngines()
	if err != nil {
		return err
	}
	sew.PopulateList(se)
	return nil
}

func (sew *SecretEngineView) PopulateList(se []string) {
	for _, engine := range se {
		selected := func() {
			sew.tui.ShowSecretsView(engine)
		}
		sew.list.Add(engine, "", selected)
	}
}
