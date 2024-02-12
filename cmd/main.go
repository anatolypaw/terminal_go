// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"image/color"
	"log"
	"sgui"
	"sgui/widget"
)

func main() {
	ui, err := sgui.New()
	if err != nil {
		log.Panic(err)
	}

	button := widget.NewButton(50, 50, "click me", nil)
	button2 := widget.NewButton(100, 50, "click me", nil)

	ind := widget.NewIndicator(30)
	ind.AddState(color.RGBA{0, 255, 255, 255})

	ind2 := widget.NewIndicator(30)
	ind2.AddState(color.RGBA{0, 255, 0, 255})

	ui.AddWidget(400, 40, button)
	ui.AddWidget(400, 100, button2)

	ui.AddWidget(465, 50, ind)
	ui.AddWidget(200, 100, ind2)

	ui.Render()

	_ = ui
	_ = button

	fmt.Printf("%#v\n", ui.DisplaySize())

}
