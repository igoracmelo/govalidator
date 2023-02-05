package main

import (
	"estudiosol/handler"
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/verify", handler.HandleVerify)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen and serve: %s", err.Error())
		os.Exit(1)
	}
}
