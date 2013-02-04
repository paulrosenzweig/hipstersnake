package main

import (
	"flag"
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
	portPtr := flag.Int("port", 80, "server port")
	flag.Parse()

	http.HandleFunc("/", mainHandler)
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.Handle("/ws", connection.ConnectionHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *portPtr), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
