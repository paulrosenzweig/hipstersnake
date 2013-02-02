package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"snake/connection"
)

func mainHandler(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("index.html")
	if err == nil {
		io.Copy(w, file)
	}
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.Handle("/ws", connection.ConnectionHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", 80), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}