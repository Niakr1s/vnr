package tray

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
)

func Run(tooltip string) {
	fullTooltip := fmt.Sprintf("Visual Novel Reader - %s", tooltip)

	go func() {
		systray.Run(func() {
			systray.SetTitle("VNR")
			systray.SetTooltip(fullTooltip)
			quit := systray.AddMenuItem("Quit", "Quit the app")

			go func() {
				<-quit.ClickedCh
				systray.Quit()
			}()
		}, func() {
			os.Exit(0)
		})
	}()
}
