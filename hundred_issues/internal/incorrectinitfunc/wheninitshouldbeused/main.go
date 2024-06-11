package main

import "database/sql"

/*
WRONG WAY
*/
// var db *sql.DB

// func init() {
// 	dataSourceName := os.Getenv("MYSQL_DATA_SOURCE_NAME")

// 	d, err := sql.Open("mysql", dataSourceName)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	err = d.Ping()
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	db = d
// }

func createClient(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
