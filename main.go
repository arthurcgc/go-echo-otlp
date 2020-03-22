package main

import (
	"log"
	"os"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	srv := http.Server{
		Handler: router,
		Addr:    ":8888",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		hostname, _ := os.Hostname()
		str := "You've hit " + hostname + "\n"
		w.Write([]byte(str))
	})

	log.Fatal(srv.ListenAndServe())
}
