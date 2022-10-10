package canon

import (
	"database/sql"
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
}

var createFormat = `
CREATE TABLE IF NOT EXISTS public."%s" (
    tm timestamp with time zone NOT NULL,
    p1 bigint NOT NULL,
    p2 real NOT NULL,
    p3 boolean NOT NULL
);
`
var insertFormat = `
INSERT INTO public."%s" (tm,p1,p2,p3) VALUES ('%s',%d,%f,'%v');
`
var readFormat = `
SELECT avg(p1),avg(p2) FROM public."%s" WHERE tm BETWEEN '%s' AND '%s'; 
`

func (d *datarecord) newData() {
	d.tm = time.Now()
	d.p1 = rand.Int()
	d.p2 = rand.Float32() * float32(d.p1)
	d.p3 = rand.Int()%2 == 0
}
func (t *Table) createifneed() {
	create := fmt.Sprintf(createFormat, t.name)
	t.db.Exec(create)
}
func (t *Table) newData() {
	t.data.newData()
	str := fmt.Sprintf(insertFormat, t.name,
		string(pq.FormatTimestamp(t.data.tm)),
		t.data.p1, t.data.p2, t.data.p3)
	// logger.Debug.Print(str)
	_, err := t.db.Exec(str)
	if err != nil {
		logger.Error.Print(err.Error())
	}
}
func (t *Table) readData() []datarecord {
	res := make([]datarecord, 0)
	str := fmt.Sprintf(readFormat, t.name, string(pq.FormatTimestamp(time.Now().Add(-time.Hour))),
		string(pq.FormatTimestamp(time.Now())))
	rows, err := t.db.Query(str)
	if err != nil {
		logger.Error.Print(err.Error())
		return res
	}
	var d datarecord
	for rows.Next() {
		rows.Scan(&d.p1, &d.p2)
		res = append(res, d)
	}
	return res
}
