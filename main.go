package main

import "net/http"

func main() {
	serverMux := http.NewServeMux()
	s := http.Server{
		Handler: serverMux,
		Addr:    ":8080",
	}
	s.ListenAndServe()
}
