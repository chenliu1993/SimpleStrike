package main

import (
	"flag"
	"os"

	"github.com/chenliu1993/SimpleStrike/utils"
)

func main() {
	client := flag.Int("LocalPort", utils.ClientPort, "client access port")
	server := flag.Int("RemotePort", utils.ServerPort, "Service access port")
	flag.Parse()
	if flag.NFlag() != 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	ts := utils.NewServer()
	ts.StartService(*client, *server)
	select {}
}
