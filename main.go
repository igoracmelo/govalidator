package main

import (
	"estudiosol/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/verify", handler.HandleVerify)
	http.ListenAndServe(":8080", nil)
}
