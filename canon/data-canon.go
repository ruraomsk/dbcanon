package canon

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"
	"github.com/ruraomsk/ag-server/logger"
)

type Table struct {
	data datarecord
	db   *sql.DB
	name string
}
type datarecord struct {
	tm time.Time
	p1 int
	p2 float32
	p3 bool
	js []Status
}

var createFormat = `
CREATE TABLE IF NOT EXISTS %s (
    tm timestamp with time zone NOT NULL primary key,
    p1 bigint NOT NULL,
    p2 real NOT NULL,
    p3 boolean NOT NULL,
	js JSONB NOT NULL DEFAULT '{}'
)
WITH (
    autovacuum_enabled = FALSE
)
TABLESPACE pg_default;
`
var insertFormat = `
INSERT INTO %s (tm,p1,p2,p3,js) VALUES ('%s',%d,%f,'%v','%s');
`
var readFormat = `
SELECT js FROM %s WHERE tm BETWEEN '%s' AND '%s'; 
`
var getSizeFormat = `
SELECT count(*) FROM %s;
`

func (d *datarecord) newData() {
	d.tm = time.Now()
	d.p1 = rand.Intn(2000)
	d.p2 = rand.Float32() * float32(d.p1)
	d.p3 = rand.Int()%2 == 0
	d.js = make([]Status, 0)
	for i := 0; i < (5 + rand.Intn(10)); i++ {
		d.js = append(d.js, Status{Tact_tick: i, Program_number: d.p1, Tact_number: d.p1 % 7})
	}
}
func (t *Table) createifneed() {
	create := fmt.Sprintf(createFormat, t.name)
	for {
		_, err := t.db.Exec(create)
		if err != nil {
			logger.Error.Print(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
}
func (t *Table) newData() {
	t.data.newData()
	js, _ := json.Marshal(t.data.js)
	str := fmt.Sprintf(insertFormat, t.name,
		string(pq.FormatTimestamp(t.data.tm)),
		t.data.p1, t.data.p2, t.data.p3, string(js))
	// logger.Debug.Print(str)
	_, err := t.db.Exec(str)
	if err != nil {
		logger.Error.Print(err.Error())
	}
}
func (t *Table) readData() ([]datarecord, error) {
	res := make([]datarecord, 0)
	str := fmt.Sprintf(readFormat, t.name,
		string(pq.FormatTimestamp(time.Now().Add(-5*time.Minute))),
		string(pq.FormatTimestamp(time.Now())))
	rows, err := t.db.Query(str)
	if err != nil {
		logger.Error.Print(err.Error())
		return res, err
	}
	var d datarecord
	var buf []byte
	for rows.Next() {
		rows.Scan(&buf)
		err := json.Unmarshal(buf, &d.js)
		if err != nil {
			logger.Error.Print(err.Error())
		} else {
			res = append(res, d)
		}
	}
	if len(res) == 0 {
		logger.Error.Printf("Table %s empty seek", t.name)
	}
	return res, nil
}
func (t *Table) getSize() (int, error) {
	str := fmt.Sprintf(getSizeFormat, t.name)
	rows, err := t.db.Query(str)
	if err != nil {
		logger.Error.Print(err.Error())
		return -1, err
	}
	var d int
	for rows.Next() {
		rows.Scan(&d)
	}
	return d, nil
}
