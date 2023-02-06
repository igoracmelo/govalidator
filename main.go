package main

import (
	"estudiosol/handler"
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/verify", handler.HandleVerify)
	fmt.Println("Listening on localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen and serve: %s\n", err.Error())
		os.Exit(1)
	}
}
