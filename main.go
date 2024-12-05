package main

import (
	"vaultview/pkg/tui"
)

func main() {
	tui := tui.NewTui()
	tui.Init()
	if err := tui.Run(); err != nil {
		panic(err)
	}
}
