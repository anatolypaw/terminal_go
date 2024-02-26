// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"os"
	"terminal/internal/app"
	"terminal/internal/guiview"
)

func main() {
	exit := make(chan os.Signal, 1)

	app := app.New()
	go app.Run()

	gui := guiview.New(&app)
	go gui.Run()

	<-exit
}
