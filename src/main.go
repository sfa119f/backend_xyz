package main

import (
	"fmt"
	"net/http"

	"github.com/sfa119f/backend_xyz/src/database"

	"github.com/gorilla/mux"
)

func main() {
	// init db
	database.InitDB()

	// init router
	router := mux.NewRouter()

	// routes
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	}).Methods(http.MethodGet)

	http.ListenAndServe(":8000", router)
}