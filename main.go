package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	repo := os.Getenv("GIT_IAC_REPO")
	if repo == "" {
		log.Fatal("env 'GIT_IAC_REPO' not set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "func-deploy")
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "pong")
	})
	http.HandleFunc("/coding", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "success")
	})

	port := os.Getenv("FC_SERVER_PORT")
	if port == "" {
		port = "9000"
	}

	log.Println("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
