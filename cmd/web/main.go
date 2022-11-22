package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP port address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("../../ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// mux.HandleFunc("/", home)
	mux.Handle("/", http.HandlerFunc(home))
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	infoLog.Printf("Starting server on %s", *addr)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	// err := http.ListenAndServe(*addr, mux)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
