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
	host := flag.String("host", "127.0.0.1", "server addr")
	remotePort := flag.Int("RemotePort", utils.RemotePort, "transfer server addr")
	localPort := flag.Int("LocalPort", utils.LocalPort, "local server addr")
	flag.Parse()
	if flag.NFlag() != 3 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	remote, err := net.Dial("tcp", fmt.Sprintf("%s:%d", utils.Host, remotePort))
	if err != nil {
		log.Fatal(err)
		return
	}
	go utils.ClientHandler(remote, *localPort)
	select {}
}
