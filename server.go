package main

import (
	"flag"
	"os"

	"github.com/chenliu1993/SimpleStrike/utils"
)

func main() {
	localPort := flag.Int("LocalPort", utils.ClientPort, "client access port")
	remotePort := flag.Int("RemotePort", utils.ServerPort, "Service access port")
	flag.Parse()
	if flag.NFlag() != 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	ts := utils.NewServer()
	ts.StartService(*localPort, *remotePort)
	select {}
}
