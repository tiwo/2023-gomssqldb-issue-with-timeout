package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/microsoft/go-mssqldb"
	mssql "github.com/microsoft/go-mssqldb"
)

func opendatabase() *sql.DB {

	mssql.SetLogger(log.Default())

	connString := os.Getenv("CONNECTIONSTRING")

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Printf("Error from sql.Open: %#v", err.Error())
	}

	db.SetConnMaxLifetime(21 * time.Second)
	return db
}

func query_with_timeout(db *sql.DB, query string) (string, bool) {
	log.Printf("query_with_timeout(%#v)", query)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// the context should time out after five seconds;
	// but for good measure, we'll also cancel it after 12 seconds:
	timer := time.AfterFunc(12*time.Second, func() {
		log.Printf("Cancelling query... state of ctx is: %#v)", ctx.Err())
		cancel()
	})
	defer timer.Stop()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error from db.Query: %#v", err.Error())
		return "", false
	}
	defer rows.Close()

	if !rows.Next() {
		log.Printf("Zero rows in the result set.")
		return "", false
	}

	var result string
	err = rows.Scan(&result)
	if err != nil {
		log.Printf("Error from rows.Scan: %#v", err.Error())
		return "", false
	}

	log.Printf("query_with_timeout returning %#v", result)
	return result, true
}

func main() {
	log.SetFlags(log.Ltime | log.LUTC | log.Lshortfile)
	log.Printf("main()")
	defer log.Printf("main() done")

	catch := make(chan os.Signal, 1)
	signal.Notify(catch, os.Interrupt)
	go func() {
		<-catch
		log.Printf("Received SIGINT; exiting.")
		os.Exit(1)
	}()

	db := opendatabase()
	defer db.Close()

	// After 50 seconds, close the database:
	timer := time.AfterFunc(50*time.Second, func() {
		log.Printf("Closing database")
		db.Close()
	})
	defer timer.Stop()

	query_with_timeout(db, "SELECT 'first';")

	log.Printf("to reproduce, disconnect from your network now")
	time.Sleep(5 * time.Second)

	query_with_timeout(db, "SELECT 'second';")

}
