package main

import (
	"gin_exercise/config"
	"gin_exercise/httpserver"
	"github.com/cihub/seelog"
	"os"
	"os/signal"
	"syscall"
)

func Init() {
	config.Initlog()
	config.Initdatabase()
	httpserver.Initroute()
}

func main() {
	Init()
	defer seelog.Flush()

	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)

	// soft kill
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	signal.Stop(signalChan)
	seelog.Infof("See you next time at %d !", config.GlobalConfig.ListenPort)

}
