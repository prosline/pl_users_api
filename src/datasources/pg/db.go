package pg

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	pg_dsn               = "pg_dsn"
	pg_dsn_user_name     = "pg_dsn_user_name"
	pg_dsn_user_password = "pg_dsn_user_password"
	pg_dsn_host          = "pg_dsn_host"
	pg_dsn_schema        = "pg_dsn_schema"

//DSN = "postgres://marciodasilva@127.0.0.1:5432/plapi?sslmode=disable"
// pg_dsn="postgres://%s@%s:5432/%s?sslmode=disable"
/*export pg_dsn="postgres://%s@%s:5432/%s?sslmode=disable"
export pg_dsn_user_name="marciodasilva"
export pg_dsn_user_password=""
export pg_dsn_host="127.0.0.1"
export pg_dsn_schema="plapi"
*/
)

var (
	err      error
	ClientDB *sqlx.DB
	pgDsn    = os.Getenv(pg_dsn)
	userName = os.Getenv(pg_dsn_user_name)
	//usePassward = os.Getenv(pg_dsn_user_password)
	userHost   = os.Getenv(pg_dsn_host)
	userSchema = os.Getenv(pg_dsn_schema)
)

func init() {
	fmt.Println(pgDsn)
	fmt.Println(userName)
	fmt.Println(userHost)
	fmt.Println(userSchema)

	ClientDBConnection := fmt.Sprintf("postgres://%v@%v:5432/%v?sslmode=disable", userName, userHost, userSchema)

	ClientDB, err = sqlx.Connect("postgres", ClientDBConnection)
	if err != nil {
		panic(err)
	}
	if errPing := ClientDB.Ping(); errPing != nil {
		panic(errPing)
	}

	log.Printf("Successful Connection to %s \n", ClientDB.DriverName())
}
