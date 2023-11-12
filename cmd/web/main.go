package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type applicatoin struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":3000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &applicatoin{
		infoLog:  infoLog,
		errorLog: errorLog,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("starting server on port %s", *addr)
	err := srv.ListenAndServe()

	errorLog.Fatal(err)

}
