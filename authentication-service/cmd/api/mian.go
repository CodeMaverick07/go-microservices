package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)
 
const webPort = 80
var count int

type Config struct{
	DB *sql.DB
	Models data.Models
}

func main() {

	conn := connectDB()

	if conn == nil {
		log.Panic("could not connect to the database")
	}

	app:= Config{
		DB: conn,
		Models: data.New(conn),
	}
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
	panic(err)
	}
}

func openDB(dsn string) (*sql.DB,error) {
  db,err:=sql.Open("pgx", dsn)

  if err !=nil {
	return nil, err
  }
  err= db.Ping()
  if err != nil {
	return nil, err
  } 
  return db,nil
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")
	fmt.Print(dsn)

	for {
		connect ,err := openDB(dsn)
		if err != nil {
			log.Println("progress is not ready yet")
			count++
		} else {
			log.Println("connected to the database")
			return connect
		}
		if count > 10 {
		  
		  return nil
		}
		log.Println("backing off for 2 sec")
		time.Sleep(2 * time.Second)
		continue 
	}
	
}