package canon

import (
	"database/sql"
	"fmt"
	"time"

	//Инициализатор постргресса
	_ "github.com/lib/pq"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/dbcanon/setup"
)

func Cannon(name string) {
	var err error
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		setup.Set.DataBase.Host, setup.Set.DataBase.User,
		setup.Set.DataBase.Password, setup.Set.DataBase.DBname)
	table := new(Table)
	table.name = name
	for {
		table.db, err = sql.Open("postgres", dbinfo)
		if err != nil {
			logger.Error.Print(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		err = table.db.Ping()
		if err != nil {
			logger.Error.Print(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		table.createifneed()
		step := time.NewTicker(time.Duration(setup.Set.Step) * time.Millisecond)

		for {
			<-step.C
			table.newData()

		}
	}
}
func Reader(name string, interval int) {
	var err error
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		setup.Set.DataBase.Host, setup.Set.DataBase.User,
		setup.Set.DataBase.Password, setup.Set.DataBase.DBname)
	table := new(Table)
	table.name = name
	for {
		table.db, err = sql.Open("postgres", dbinfo)
		if err != nil {
			logger.Error.Print(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		err = table.db.Ping()
		if err != nil {
			logger.Error.Print(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		step := time.NewTicker(time.Duration(interval) * time.Millisecond)
		for {
			<-step.C
			table.readData()

		}
	}

}
