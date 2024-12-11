package tui

import (
	"vaultview/pkg/constants"
	"vaultview/pkg/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Info struct {
	*tview.Table
	tui *Tui
}

func NewInfo(tui *Tui) *Info {
	info := &Info{
		Table: tview.NewTable(),
		tui:   tui,
	}
	info.layout()
	return info
}

func (it *Info) layout() {
	for row, info := range []string{"VaultView Rev:", "Vault Rev:", "Vault Addr:", "Sealed:", "Token Policies:", "Token Expires:"} {
		it.Table.SetCell(row, 0, it.getInfoCell(info))
		it.Table.SetCell(row, 1, it.getInfoValueCell(constants.NAValue))
	}
}

func (it *Info) getInfoCell(info string) *tview.TableCell {
	cell := tview.NewTableCell(info)
	cell.SetTextColor(tcell.ColorLime)
	cell.SetAlign(tview.AlignLeft)
	return cell
}

func (it *Info) getInfoValueCell(info string) *tview.TableCell {
	cell := tview.NewTableCell(info)
	cell.SetAlign(tview.AlignLeft)
	return cell
}

func (it *Info) setCell(row int, newValue string) int {
	if newValue != "" {
		it.GetCell(row, 1).SetText(newValue)
	}
	return row + 1
}

func (it *Info) UpdateInfoTable(data models.Info) {
	it.tui.QueueUpdateDraw(func() {
		nextRow := it.setCell(0, data.VaultViewRev)
		nextRow = it.setCell(nextRow, data.VaultRev)
		nextRow = it.setCell(nextRow, data.VaultAddr)
		nextRow = it.setCell(nextRow, data.Sealed)
		nextRow = it.setCell(nextRow, data.TokenPolicies)
		nextRow = it.setCell(nextRow, parseTime(data.TokenExpirationTime))
	})
}
