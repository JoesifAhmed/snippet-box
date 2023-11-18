package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/JoesifAhmed/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type applicatoin struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":3000", "HTTP network address")

	dsn := flag.String("dsn", "web:MySQLWeb&911@/snippetbox?parseTime=true", "MySQL Database")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTamplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &applicatoin{
		infoLog:       infoLog,
		errorLog:      errorLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("starting server on port %s", *addr)
	err = srv.ListenAndServe()

	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
