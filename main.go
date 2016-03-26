package main

import (
	"flag"
	"github.com/dinever/dingo/app"
	"github.com/dinever/dingo/cmd"
)

func main() {
	if !cmd.CheckInstall() {
		cmd.Install()
	}
	portPtr := flag.String("port", "8000", "The port number for Dingo to listen to.")
	flag.Parse()
	Dingo.Init()
	Dingo.Run(*portPtr)
}
