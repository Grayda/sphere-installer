package main

import (
	"fmt"
	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/go-ninja/support"
)

var host = config.String("localhost", "led.host")
var port = config.Int(3115, "led.remote.port")

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
