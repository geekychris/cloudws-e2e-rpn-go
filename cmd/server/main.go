// Command server runs the RPN scientific calculator HTTP service.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/geekychris/cloudws-e2e-rpn-go/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/calculate", api.CalculateHandler)

	addr := ":" + port
	log.Printf("RPN calculator service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
