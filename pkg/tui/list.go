package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

type List struct {
	list        *tview.List
	showSecText bool
}

func NewList(title string, tui *Tui) *List {
	li := &List{
		list:        list(title),
		showSecText: false,
	}

	return li
}

func list(title string) *tview.List {
	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetTitle(fmt.Sprintf(" %v ", title)).SetBorder(true)
	return list
}

func (l *List) Clear() {
	l.list.Clear()
}

func (l *List) EnableSecText() {
	l.showSecText = !l.showSecText
	l.list.ShowSecondaryText(l.showSecText)
}

func (l *List) SetTitle(title string) {
	l.list.SetTitle(fmt.Sprintf(" %v ", title))
}

func (l *List) Hydrate(items []string, selected ...string) {
	l.list.Clear().SetOffset(0, 0)
	selectedItem := -1
	for i, name := range items {
		if len(selected) > 0 {
			if selected[0] == name {
				selectedItem = i
			}
		}
		l.list.AddItem(name, "", 0, nil)
	}
	if selectedItem >= 0 {
		l.list.SetCurrentItem(selectedItem)
	}
}

func (l *List) Add(item, secItem string, f func()) {
	l.list.AddItem(item, secItem, 0, f)
}

func (l *List) List() *tview.List {
	return l.list
}

func (l *List) NextItem() {
	nextItem := l.list.GetCurrentItem() + 1
	if nextItem >= l.list.GetItemCount() {
		l.list.SetCurrentItem(0)
	} else {
		l.list.SetCurrentItem(nextItem)
	}
}

func (l *List) getItemText() string {
	i := l.list.GetCurrentItem()
	main, _ := l.list.GetItemText(i)
	return main
}
