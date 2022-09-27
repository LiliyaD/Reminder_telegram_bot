package postgres

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	Host = "localhost"
	Port = 5432

	MaxConnIdleTime = time.Minute
	MaxConnLifetime = time.Hour
	MinConns        = 2
	MaxConns        = 4
)

var (
	User, Password, DBname string

	fieldsForOrder map[string]string

	_testing = false
)

func init() {
	if _testing {
		return
	}

	var exists bool

	User, exists = os.LookupEnv("POSTGRESQL_USER")
	if !exists {
		log.Fatal(errors.New("Environment variable POSTGRESQL_USER doesn't exist"))
	}

	Password, exists = os.LookupEnv("POSTGRESQL_PASS")
	if !exists {
		log.Fatal(errors.New("Environment variable POSTGRESQL_PASS doesn't exist"))
	}

	DBname, exists = os.LookupEnv("POSTGRESQL_DB")
	if !exists {
		log.Fatal(errors.New("Environment variable POSTGRESQL_DB doesn't exist"))
	}

	fieldsForOrder = make(map[string]string)
	fieldsForOrder["name"] = "act_name"
	fieldsForOrder["beginDate"] = "begin_date"
	fieldsForOrder["endDate"] = "end_date"
	fieldsForOrder["timesPerDay"] = "times_per_day"
	fieldsForOrder["quantityPerTime"] = "quantity_per_time"
}
