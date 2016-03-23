package main

import (
	"github.com/dinever/dingo/app"
	"github.com/dinever/dingo/cmd"
)

func main() {
	if !cmd.CheckInstall() {
		cmd.Install()
	}
	Dingo.Init()
	Dingo.Run()
}
