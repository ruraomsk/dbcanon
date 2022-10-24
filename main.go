package main

import (
	"embed"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/ruraomsk/ag-server/logger"

	"github.com/ruraomsk/dbcanon/canon"
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
	setup.Set.Work = true
	go canon.Statistics()
	fmt.Println("Готова статистика....")
	time.Sleep(2 * time.Second)
	fmt.Println("Начали писать....")
	for i := 0; i < setup.Set.TablesCount-1; i++ {
		go canon.Cannon(fmt.Sprintf("table%d", i), true)
	}
	//Эта для проверки спинлока
	go canon.Cannon(fmt.Sprintf("table%d", setup.Set.TablesCount-1), false)
	time.Sleep(20 * time.Second)
	fmt.Println("Начали читать....")
	for i := 0; i < setup.Set.TablesAvg; i++ {
		j := rand.Intn(setup.Set.Maximum)
		if j <= 0 {
			j = setup.Set.Maximum
		}
		go canon.Reader(fmt.Sprintf("table%d", i), j)
	}
	//Эта для проверки спинлока
	go canon.Reader(fmt.Sprintf("table%d", setup.Set.TablesCount-1), -1)

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
		setup.Set.Work = false
		time.Sleep(5 * time.Second)
		break loop
	}
	logger.Info.Print("Stop dbcanon")
}
