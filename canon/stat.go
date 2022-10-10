package canon

import (
	"fmt"
	"time"
)

type stat struct {
	name string
	op   string
	long time.Duration
}
type extstat struct {
	name   string
	wcount int64
	wsumm  int64
	rcount int64
	rsumm  int64
}

var cstat chan stat
var maps map[string]*extstat

func Statistics() {
	cstat = make(chan stat, 1000)
	maps = make(map[string]*extstat)
	oneMinute := time.NewTicker(time.Minute)
	for {
		select {
		case <-oneMinute.C:
			var (
				rc int64
				rs int64
				wc int64
				ws int64
			)
			for _, v := range maps {
				rc += v.rcount
				rs += v.rsumm
				wc += v.wcount
				ws += v.wsumm
				v.rcount = 0
				v.rsumm = 0
				v.wcount = 0
				v.wsumm = 0
			}
			if rc != 0 {
				fmt.Printf("\t%fs", (float64(rs)/float64(rc))/1000000000)
			} else {
				fmt.Printf("\tnot\t")
			}
			if wc != 0 {
				fmt.Printf("\t%fs", (float64(ws)/float64(wc))/1000000000)
			} else {
				fmt.Printf("\tnot")
			}
			fmt.Print("\n")

		case s := <-cstat:
			var es *extstat
			var is bool
			es, is = maps[s.name]
			if !is {
				es = &extstat{name: s.name, wcount: 0, wsumm: 0, rcount: 0, rsumm: 0}
			}
			switch s.op {
			case "read":
				es.rcount++
				es.rsumm += int64(s.long)
			case "write":
				es.wcount++
				es.wsumm += int64(s.long)
			}
			maps[es.name] = es
		}
	}
}
