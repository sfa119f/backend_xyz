package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sfa119f/backend_xyz/src/database"
	// "github.com/sfa119f/backend_xyz/src/dictionary"
	"github.com/sfa119f/backend_xyz/src/handler"
	"github.com/sfa119f/backend_xyz/src/utils"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// init db
	database.InitDB()

	// init router
	router := mux.NewRouter()

	// route without middleware
	r_wm := router.PathPrefix("/api").Subrouter()
	r_wm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	}).Methods(http.MethodGet)
	r_wm.HandleFunc("/customer", handler.InsertCustomer).Methods(http.MethodPost)
	r_wm.HandleFunc("/login", handler.Login).Methods(http.MethodPost)

	// route using middleware
	r_um := router.PathPrefix("/api").Subrouter()
	r_um.Use(utils.MiddlewareJWTAuthorization)
	r_um.HandleFunc("/customer", handler.UpdateCstExceptPass).Methods(http.MethodPut)
	r_um.HandleFunc("/customer/password", handler.UpdateCstPass).Methods(http.MethodPut)
	r_um.HandleFunc("/customer/details", handler.InsertUpdateCstDetails).Methods(http.MethodPost)
	r_um.HandleFunc("/tenorLimit", handler.GetTenorByIdCust).Methods(http.MethodGet)
	r_um.HandleFunc("/transaction", handler.InsertTransaction).Methods(http.MethodPost)
	
	port :=  ":" + os.Getenv("XYZ_PORT")
	server := new(http.Server)
	server.Handler = router
	server.Addr = port

	fmt.Println("Starting server at", server.Addr)
	server.ListenAndServe()
}