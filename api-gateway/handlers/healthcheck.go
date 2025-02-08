package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Printf("Successful Health Check!")
	fmt.Fprintln(w, "Successful Health Check!")
}
