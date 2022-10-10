package setup

var (
	Set *Setup
)

type Setup struct {
	LogPath     string   `toml:"logpath"`
	DataBase    DataBase `toml:"dataBase"`
	TablesCount int      `toml:"count"` //Кол-во таблиц для записи
	Step        int      `toml:"step"`  //Интервал записи новых значений в таблицы (миллисекунты)
	Maximum     int      `toml:"max"`   //Максимальный интервал порождения запроса на агрегацию по времени
}

// DataBase настройки базы данных postresql
type DataBase struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBname   string `toml:"dbname"`
	Step     int    `toml:"step"`
}

func init() {
}
