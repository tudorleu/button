package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
	"log"
)

var DbMap *gorp.DbMap

func InitDb(dbUrl string, dropTables bool) {
	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalln(err)
	}

	DbMap = &gorp.DbMap{Db: dbConn, Dialect: gorp.PostgresDialect{}}
	//	DbMap.TraceOn("[gorp]", log.New(os.Stdout, "db: ", log.Ltime))

	DbMap.AddTableWithName(User{}, "users").SetKeys(true, "id")
	DbMap.AddTableWithName(Transfer{}, "transfers").SetKeys(true, "id")

	if dropTables {
		err = DbMap.DropTables()
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = DbMap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln("Error while attempting to create tables")
	}
}

func CloseDb() {
	DbMap.Db.Close()
}
