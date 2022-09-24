package tray

import (
	"os"

	"github.com/getlantern/systray"
)

func Run() {
	go func() {
		systray.Run(onReady(), onExit)
	}()
}

func onReady() func() {
	return func() {
		systray.SetTitle("VNR")
		systray.SetTooltip("Visual Novel Reader")
		quit := systray.AddMenuItem("Quit", "Quit the app")

		go func() {
			<-quit.ClickedCh
			systray.Quit()
		}()
	}
}

func onExit() {
	os.Exit(0)
}
