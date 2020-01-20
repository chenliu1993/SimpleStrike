package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/chenliu1993/SimpleStrike/utils"
)

func main() {
	Host := flag.String("host", utils.Host, "server addr")
	remotePort := flag.Int("RemotePort", utils.RemotePort, "transfer server addr")
	localPort := flag.Int("LocalPort", utils.LocalPort, "local server addr")
	flag.Parse()
	if flag.NFlag() != 3 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	remote, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *remotePort))
	if err != nil {
		log.Fatal(err)
		return
	}
	go utils.ClientHandler(remote, *localPort)
	select {}
}
