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

func Cannon(name string, isStat bool) {
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
		// time.Sleep(time.Duration((i+1)*100) * time.Millisecond)
		// step := time.NewTicker(time.Duration(setup.Set.Step) * time.Millisecond)

		for {
			time.Sleep(time.Duration(setup.Set.Step) * time.Millisecond)
			if !setup.Set.Work {
				table.db.Close()
				return
			}
			st := time.Now()
			table.newData()
			if isStat {
				cstat <- stat{name: table.name, op: "write", long: time.Now().Sub(st)}
			}

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
		var step *time.Ticker
		if interval < 0 {
			step = time.NewTicker(time.Minute)
		} else {
			step = time.NewTicker(time.Duration(interval) * time.Minute)
		}
		for {
			<-step.C
			if !setup.Set.Work {
				table.db.Close()
				return
			}
			if interval > 0 {
				st := time.Now()
				_, err = table.readData()
				if err == nil {
					cstat <- stat{name: table.name, op: "read", long: time.Now().Sub(st)}
				}
			} else {
				//Ловим спинлок
				count, err := table.getSize()
				if err != nil {
					return
				}
				logger.Info.Printf("Table %s counted %d", table.name, count)

			}
		}
	}
}
