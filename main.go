package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/ruraomsk/ag-server/logger"

	"github.com/ruraomsk/dbcanon/setup"
)

var (
	//go:embed config
	config embed.FS
)

func init() {
	setup.Set = new(setup.Setup)
	if _, err := toml.DecodeFS(config, "config/config.toml", &setup.Set); err != nil {
		fmt.Println("Dissmis config.toml")
		os.Exit(-1)
		return
	}
	os.MkdirAll(setup.Set.LogPath, 0777)
	if err := logger.Init(setup.Set.LogPath); err != nil {
		log.Panic("Error logger system", err.Error())
		return
	}
}

func main() {
	logger.Info.Print("Start dbcanon")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP)
loop:
	for {
		<-c
		fmt.Println("\nWait make abort...")
		// uptransport.DebugStopAmi <- 1
		// time.Sleep(time.Second)
		time.Sleep(5 * time.Second)
		break loop
	}
	logger.Info.Print("Stop dbcanon")
}