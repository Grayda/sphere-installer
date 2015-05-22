package main

import (
	"fmt"
	"github.com/ninjasphere/go-ninja/support"
)

func main() {

	app := &App{}
	err := app.Init(info)
	if err != nil {
		fmt.Println("failed to initialize app:", err)
	}

	err = app.Export(app)
	if err != nil {
		fmt.Println("failed to export app:", err)
	}

	support.WaitUntilSignal()
}
