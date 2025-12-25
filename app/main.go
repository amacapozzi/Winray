//go:build windows

package main

import (
	"winray-app/internal/hotkey"
	"winray-app/internal/index"
	"winray-app/internal/ui"
)

func main() {
	go index.BuildInitial()

	uiInstance := ui.New()

	hotkey.Loop(uiInstance.Toggle)

	select {}
}
