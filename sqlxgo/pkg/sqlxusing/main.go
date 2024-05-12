package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	db1, err := creation()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", db1.Stats())

	fmt.Println()
	db2, err := openAndConnect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", db2.Stats())

	fmt.Println()
	if err := dropTable(db1); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	if err := createTable(db1); err != nil {
		log.Fatal(err)
	}

	insertions(db1)

	fmt.Println()
	if err := queryxUsing1(db1); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	if err := getAndSelectUsing(db1); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	if err := inUsing(db1); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	if err := namedQueryUsing(db1); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	if err := sliceMapScan(db1); err != nil {
		log.Fatal(err)
	}

}

func creation() (*sqlx.DB, error) {
	const connCredentials = "postgres://admin:qwerty12345@localhost:5432/dbpw"
	db, err := sqlx.Open("pgx", connCredentials)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func openAndConnect() (*sqlx.DB, error) {
	const connCredentials = "postgres://admin:qwerty12345@localhost:5432/dbpw"

	db, err := sqlx.Connect("pgx", connCredentials)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func dropTable(db *sqlx.DB) error {
	const (
		query = `
		DROP TABLE IF EXISTS place;
		`
	)

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

func createTable(db *sqlx.DB) error {
	const (
		query = `
		CREATE TABLE IF NOT EXISTS place (
			country text NOT NULL,
			city TEXT,
			telcode INTEGER,

			PRIMARY KEY (country)
		)
		`
	)

	// _, err := db.Exec(query)
	// if err != nil {
	// 	return err
	// }
	db.MustExec(query)

	return nil
}

func insertions(db *sqlx.DB) {
	const (
		query1 = `
		INSERT INTO place(country, telcode)
		VALUES
			($1, $2)
		;
		`

		query2 = `
		INSERT INTO place(country, city, telcode)
		VALUES
			($1, $2, $3)
		;
		`
	)

	db.MustExec(query1, "Hong Kong", 852)
	db.MustExec(query1, "Singapore", 65)
	db.MustExec(query2, "South Africa", "Johannestburg", 27)
}

type Place struct {
	Country       string
	City          sql.NullString
	TelephoneCode sql.NullInt32 `db:"telcode"`
}

func queryxUsing1(db *sqlx.DB) error {
	const (
		query = `
		SELECT
			*
		FROM
			place
		`
	)

	var (
		curPlace Place
	)

	rows, err := db.Queryx(query)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		if err := rows.StructScan(&curPlace); err != nil {
			return err
		}
		fmt.Printf("%v\n", curPlace)
	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func getAndSelectUsing(db *sqlx.DB) error {
	var (
		place  Place
		places []Place

		count int
		names []string
	)

	if err := db.Get(&place, `SELECT * FROM place LIMIT 1`); err != nil {
		return err
	}
	fmt.Printf("%v\n", place)

	if err := db.Select(&places, `SELECT * FROM place WHERE telcode > $1`, 50); err != nil {
		return err
	}
	fmt.Printf("%v\n", places)

	if err := db.Get(&count, `SELECT count(*) FROM place`); err != nil {
		return err
	}
	fmt.Println(count)

	if err := db.Select(&names, `SELECT country FROM place`); err != nil {
		return err
	}
	fmt.Println(names)

	return nil
}

func inUsing(db *sqlx.DB) error {
	const (
		query = `
		SELECT
			*
		FROM
			place
		WHERE
			telcode IN (?)
		`
	)

	var (
		telCodes = []int{852}

		place Place
	)

	// Используем sqlx.In для подготовки запроса с учетом списка значений
	q, args, err := sqlx.In(query, telCodes)
	if err != nil {
		return err
	}

	q = db.Rebind(q)
	rows, err := db.Queryx(q, args...)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		if err := rows.StructScan(&place); err != nil {
			return err
		}

		fmt.Println(place)
	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func namedQueryUsing(db *sqlx.DB) error {
	var (
		place = Place{Country: "South Africa"}

		cities = map[string]interface{}{"city": "Johannestburg"}
	)

	rows, err := db.NamedQuery(`SELECT * FROM place WHERE country like :country`, place)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if cols, err := rows.ColumnTypes(); err != nil {
		return err
	} else {
		fmt.Println(len(cols))
	}

	rows, err = db.NamedQuery(`SELECT * FROM place WHERE city like :city`, cities)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if cols, err := rows.ColumnTypes(); err != nil {
		return err
	} else {
		fmt.Println(len(cols))
	}

	var (
		newPlace = Place{TelephoneCode: sql.NullInt32{Int32: 50, Valid: true}}
		places   []Place

		nstmt *sqlx.NamedStmt
	)

	nstmt, err = db.PrepareNamed(`SELECT * FROM place WHERE telcode > :telcode`)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	if err = nstmt.Select(&places, newPlace); err != nil {
		return err
	}
	fmt.Println(places)

	return nil
}

func sliceMapScan(db *sqlx.DB) error {
	rows, err := db.Queryx(`SELECT * FROM place`)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		if cols, err := rows.SliceScan(); err != nil {
			return err
		} else {
			fmt.Println(cols)
		}

	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return err
	}

	var (
		results = make(map[string]interface{})
	)
	rows, err = db.Queryx(`SELECT * FROM place`)
	for rows.Next() {
		if err := rows.MapScan(results); err != nil {
			return err
		}
	}
	fmt.Println(results)

	return nil
}
