package tui

import (
	"os"
	"time"
	"vaultview/pkg/config"
	"vaultview/pkg/constants"
	"vaultview/pkg/vault"

	"github.com/rivo/tview"
)

type Tui struct {
	App              *tview.Application
	pages            *tview.Pages
	vaultConfigModal *ModalInput
	views            map[string]View
	cfg              *config.Config
	vault            vault.VaultSvc
	main             *tview.Flex
}

func NewTui() *Tui {
	tui := &Tui{
		App:   tview.NewApplication(),
		pages: tview.NewPages(),
		cfg:   config.NewConfig(),
		views: make(map[string]View),
	}

	//modal
	tui.vaultConfigModal = NewModalInput(tui)

	//header
	header := NewHeaderView(tui)

	//content
	secretEngine := NewSecretEngineView(tui)
	secretData := NewSecretDataView(tui)
	secrets := NewSecretView(tui)

	tui.pages.AddPage(constants.ViewSecretEngines, secretEngine, true, true)
	tui.pages.AddPage(constants.ViewSecrets, secrets, true, false)
	tui.pages.AddPage(constants.ViewSecretData, secretData, true, false)

	tui.views[constants.ViewHeader] = header
	tui.views[constants.ViewSecrets] = secrets
	tui.views[constants.ViewSecretData] = secretData
	tui.views[constants.ViewSecretEngines] = secretEngine

	tui.main = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 7, 0, false).
		AddItem(tui.pages, 0, 1, true)

	return tui
}

func (tui *Tui) Init() {
	if os.Getenv("VAULT_ADDR") == "" || os.Getenv("VAULT_TOKEN") == "" {
		tui.App.SetRoot(tui.vaultConfigModal, true).EnableMouse(false)
		tui.vaultConfigModal.Init()
	} else {
		tui.cfg.UpdateVaultAddr(os.Getenv("VAULT_ADDR"))
		tui.InitVault(tui.cfg.VaultAddr, os.Getenv("VAULT_TOKEN"))
		err := tui.InitMain()
		if err != nil {
			tui.ShowErrAndStop(err)
		}
	}
}

func (tui *Tui) InitMain() error {
	tui.App.SetRoot(tui.main, true).EnableMouse(false)
	err := tui.views[constants.ViewHeader].Hydrate()
	if err != nil {
		return err
	}

	err = tui.views[constants.ViewSecretEngines].Hydrate()
	if err != nil {
		return err
	}
	return nil
}

func (tui *Tui) PublishInfo(msg string) {
	tui.views[constants.ViewHeader].(HeaderViewI).Info(msg)
}

func (tui *Tui) PublishErr(msg string) {
	tui.views[constants.ViewHeader].(HeaderViewI).Err(msg)
}

func (tui *Tui) ClearStatus() {
	tui.views[constants.ViewHeader].(HeaderViewI).Reset()
}

func (tui *Tui) ShowErrAndStop(err error) {
	go func() {
		tui.PublishErr(err.Error())
		time.Sleep(time.Second * 5)
		tui.App.Stop()
	}()
}

func (tui *Tui) ShowErrAndContinue(err error) {
	tui.PublishErr(err.Error())
	go func() {
		time.Sleep(time.Second * 5)
		tui.App.QueueUpdateDraw(func() {
			tui.ClearStatus()
		})
	}()
}

func (tui *Tui) ShowInfoAndContinue(msg string) {
	tui.PublishInfo(msg)
	go func() {
		time.Sleep(time.Second * 5)
		tui.App.QueueUpdateDraw(func() {
			tui.ClearStatus()
		})
	}()
}

func (tui *Tui) TogglePage(name string) {
	for _, pName := range tui.pages.GetPageNames(true) {
		tui.pages.HidePage(pName)
	}
	tui.pages.ShowPage(name)
}

func (tui *Tui) TogglePageAndRefresh(name string) {
	//extend this to work with all interfaces (extend View interface)
	tui.views[name].(SecretViewI).SecretsHardRefresh()
	tui.TogglePage(name)
}

func (tui *Tui) ShowSecretsView(engine string) {
	tui.TogglePage(constants.ViewSecrets)
	tui.views[constants.ViewSecrets].Hydrate(engine)
}

func (tui *Tui) ShowSecretDataView(secret, engine string) {
	tui.TogglePage(constants.ViewSecretData)
	tui.views[constants.ViewSecretData].Hydrate(secret, engine)
}

func (tui *Tui) InitVault(addr, token string) {
	var err error
	tui.vault, err = vault.NewVault(addr, token)
	if err != nil {

	}
}

func (a *Tui) QueueUpdateDraw(f func()) {
	if a.App == nil {
		return
	}
	go func() {
		a.App.QueueUpdateDraw(f)
	}()
}

func (tui *Tui) Run() error {
	return tui.App.Run()
}
