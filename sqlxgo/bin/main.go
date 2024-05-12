package main

import (
	"database/sql"

	"errors"
	"fmt"
	"log"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MultiChan struct {
	err error
}

func main() {

	var (
		db  *sql.DB
		err error
	)

	if db, err = create(); err != nil {
		fmt.Println(err)
	}
	defer func() {
		db.Close()
	}()

	if err = db.Ping(); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	if err = query1(db); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	if err = query2(db); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	if err = query3(db); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	if err = insert1(db); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	if err = delete1(db); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	lostUpdateCheck(db)

	fmt.Println()
	if err = nullVals(db); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	if err = unknownData1(db); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	if err = unknownData2(db); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	if err = txsBindConn(db); err != nil {
		fmt.Println(err)
	}
}

// Data base creation
func create() (*sql.DB, error) {
	const (
		connData = "postgres://dmitriymamykin@localhost:5432/template1"
	)

	db, err := sql.Open("pgx", connData)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Ordinary querying
func query1(db *sql.DB) error {
	const (
		queryBody = `
		select
			*
		from
			products
		where
			quantity >= $1
		`
	)

	var (
		pid      int
		quantity sql.NullInt32
		price    sql.NullInt32

		err error
	)

	// Queries the database and gets the rows returned
	rows, err := db.Query(queryBody, 15)
	if err != nil {
		return err
	}
	// If the query doesn't imply the selection of all rows
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&pid, &quantity, &price); err != nil {
			return err
		}

		fmt.Println(pid, quantity, price)
	}
	// Checks the errors after scanning
	if err = rows.Err(); err != nil {
		return err
	}
	// Closes the rows twice
	rows.Close()

	return nil
}

// Prepared statement
func query2(db *sql.DB) error {
	const (
		queryBody = `
		select
			count(sub.*)
		from
			(
				select
					*
				from
					products
				where
					price > $1
			) as sub
		`
	)

	var (
		stmt      *sql.Stmt
		rowsCount int

		err error
	)

	if stmt, err = db.Prepare(queryBody); err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(250)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&rowsCount); err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}
	rows.Close()

	fmt.Println(rowsCount)

	return nil

}

// Single-Row queries
func query3(db *sql.DB) error {
	const (
		queryBody = `
		select
			count(sub.*)
		from
			(
				select
					*
				from
					products
				where
					price > $1
			) as sub
		`
	)

	var (
		count int

		stmt *sql.Stmt
		err  error
	)

	if stmt, err = db.Prepare(queryBody); err != nil {
		return err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(100).Scan(&count); err != nil {
		return err
	}

	fmt.Println(count)

	return nil
}

func insert1(db *sql.DB) error {
	const (
		queryBody = `
		INSERT INTO
			products(quantity, price)
		VALUES
			($1, $2)
		`
	)

	var (
		stmt *sql.Stmt
		err  error
	)

	if stmt, err = db.Prepare(queryBody); err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(15, 80)
	if err != nil {
		return err
	}

	// LastInsertId isn't supported by pgx driver
	// lastInsertedRowId, err := res.LastInsertId()
	// if err != nil {
	// 	return err
	// }

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("Rows count: %d\n", rowsCount)

	return nil
}

func delete1(db *sql.DB) error {
	const (
		queryBody = `
		DELETE FROM
			products
		WHERE
			pid >= $1
		`
	)

	var (
		stmt *sql.Stmt
		err  error
	)

	if stmt, err = db.Prepare(queryBody); err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(5)
	if err != nil {
		return err
	}

	if rowsAffectedCount, err := res.RowsAffected(); err != nil {
		return err
	} else {
		fmt.Printf("Rows deleted count: %d\n", rowsAffectedCount)
	}

	return nil
}

func lostUpdateCheck(db *sql.DB) {

	var (
		wg sync.WaitGroup
		to int
	)

	for i := 0; i < 2; i++ {
		wg.Add(1)

		if i == 0 {
			to = 1500
		} else {
			to = 1750
		}

		go func(i, from, to int) {
			defer wg.Done()

			for v := range tx1(db, from, to) {
				if v.err != nil {
					fmt.Printf("%d: %v\n", i, v.err)
				}
			}
		}(i, 1750, to)

	}

	wg.Wait()
}

func tx1(db *sql.DB, from, to int) <-chan MultiChan {
	var (
		signal = make(chan MultiChan)
	)

	go func() {
		defer close(signal)

		const (
			queryBody = `
			UPDATE
				products
			SET
				price = $2
			where
				price = $1
			`
		)

		var (
			tx          *sql.Tx
			txFirstStmt *sql.Stmt
			err         error
		)

		if tx, err = db.Begin(); err != nil {
			signal <- MultiChan{err: err}
			return
		}
		defer tx.Rollback()

		if txFirstStmt, err = tx.Prepare(queryBody); err != nil {
			signal <- MultiChan{err: err}
			return
		}

		if _, err := txFirstStmt.Exec(from, to); err != nil {
			signal <- MultiChan{err: err}
			return
		}
		txFirstStmt.Close()

		if err = tx.Commit(); err != nil {
			signal <- MultiChan{err: err}
			return
		}

		signal <- MultiChan{err: errors.New("value changed")}
	}()

	return signal
}

func errors1(db *sql.DB, lgr *log.Logger) error {
	const (
		queryBody = `
		SELECT
			*
		FROM
			(
				SELECT
					quantity,
					price,
					rank() OVER (PARTITION BY quantity ORDER BY price) as rnk
				FROM
					products
			) as ranked
		WHERE
			ranked.rnk = $1
		`
	)

	var (
		stmt *sql.Stmt
		err  error

		pid             int
		quantity, price sql.NullInt32
	)

	if stmt, err = db.Prepare(queryBody); err != nil {
		return err
	}

	rows, err := stmt.Query(2)
	if err != nil {
		return err
	}
	// Writes logs after closing the rows.
	defer func() {
		if err := rows.Close(); err != nil {
			lgr.Println(err)
		}
	}()

	for rows.Next() {
		if err := rows.Scan(&pid, &quantity, &price); err != nil {
			return err
		}

		fmt.Println(pid, quantity, price)
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	return nil

}

func errors2(db *sql.DB) error {
	const (
		queryBody = `
		SELECT
			*
		FROM
			(
				SELECT
					quantity,
					price,
					rank() OVER (PARTITION BY quantity ORDER BY price) as rnk
				FROM
					products
			) as ranked
		WHERE
			ranked.rnk = $1
		LIMIT $2
		`
	)

	var (
		err error

		pid             int
		quantity, price sql.NullInt32
	)

	if err = db.QueryRow(queryBody, 2, 1).Scan(&pid, &quantity, &price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("there's no row needed: %v", err)
		}

		return fmt.Errorf("querying a row: %v", err)
	}

	return nil

}

func errors3() (*sql.DB, error) {
	const (
		connCredentials = "postgres://admin:12345@localhost:5432/dbpw"
	)

	db, err := sql.Open("pgx", connCredentials)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func nullVals(db *sql.DB) error {
	const (
		queryBody = `
		SELECT
			*
		FROM
			products
		WHERE
			price = $1
		;
		`
	)

	var (
		err error

		pid             int
		quantity, price sql.NullInt32
	)

	if err = db.QueryRow(queryBody, 3).Scan(&pid, &quantity, &price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("there's no row needed: %v", err)
		}

		return fmt.Errorf("querying a row: %v", err)
	}

	fmt.Print(pid, " ")
	if quantity.Valid {
		fmt.Print(quantity, " ")
	} else {
		fmt.Printf("\nquantity isn't presented for this row\n")
	}

	if price.Valid {
		fmt.Print(quantity, " ")
	} else {
		fmt.Printf("\n price isn't presented for this row\n")
	}

	return nil
}

func unknownData1(db *sql.DB) error {
	const (
		queryBody = `
		SELECT
			*
		FROM
			products
		;
		`
	)

	var (
		err error
	)

	rows, err := db.Query(queryBody)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	types, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	for _, v := range columns {
		fmt.Print(v, " ")
	}
	fmt.Println()

	for _, v := range types {
		fmt.Print(v.DatabaseTypeName(), " ")

	}
	fmt.Println()

	return nil

}

func unknownData2(db *sql.DB) error {
	const (
		queryBody = `
		SELECT
			*
		FROM
			products
		;
		`
	)

	var (
		err  error
		vals []interface{}
	)

	rows, err := db.Query(queryBody)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	vals = make([]interface{}, len(columns))
	for i := 0; i < len(vals); i++ {
		vals[i] = new(sql.RawBytes)
	}

	for i := 0; i < len(vals) && rows.Next(); i++ {
		if err := rows.Scan(vals...); err != nil {
			return err
		}

	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("known error: %v", err)
		}

		return fmt.Errorf("bug: %v", err)
	}

	for _, v := range vals {
		fmt.Println(*v.(*sql.RawBytes))
	}

	return nil
}

func txsBindConn(db *sql.DB) error {
	const (
		queryBody = `
		SELECT
			pid
		FROM
			products
		WHERE
			price = $1
		`
	)

	var (
		pid int
	)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(queryBody)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(1500)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		if err := rows.Scan(&pid); err != nil {
			return err
		}

		// error
		tx.Query("select * from products where pid = $1", pid)
	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}