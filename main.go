package main

import (
	"./controller"
	"log"
	"net/http"
)

func main() {
	http.Handle("/public/", http.FileServer(http.Dir("./")))

	http.HandleFunc("/", controller.Index)

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("端口错误")
	}
}
