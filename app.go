package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// App.init() initializes the app's configuration and database'
func (a *App) init(config *config) {
	// TODO: Remove below debug message
	log.Println("Getting configuration")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.username,
		config.DB.password,
		config.DB.host,
		config.DB.port,
		config.DB.name,
		config.DB.charset)

	// Create database "connection" to use for life of app

	// TODO: Remove below debug messages
	log.Println("Connection to database")
	log.Printf("DSN: %s\n", dsn)
	var err error
	a.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("Failed to open DB: ", err)
	}
	if err = a.DB.Ping(); err != nil {
		log.Fatalln("Connecting to DB failed: ", err)
	}
	defer a.DB.Close()

	a.Router = mux.NewRouter()
	a.initRoutes()
}

// initRoutes() creates all the required API routes
func (a *App) initRoutes() {
	a.Router.HandleFunc("/devices", a.listDevices).Methods("GET")
}

// run() starts the API server
func (a *App) run(addr string) {
	// TODO: Remove below debug message
	log.Println("API server starting")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) listDevices(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, map[string]string{"result": "success"})
}

// writeJSONResponse() is helper that returns JSON HTTP response
func writeJSONResponse(w http.ResponseWriter, resCode int, payload interface{}) {
	p, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(resCode)
	w.Write(p)
}
