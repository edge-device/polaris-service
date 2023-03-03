package main

import (
	"database/sql"
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.username,
		config.DB.password,
		config.DB.host,
		config.DB.port,
		config.DB.name,
		config.DB.charset)

	// Create database "connection" to use for life of app
	var err error
	a.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("Failed to open DB: ", err)
	}
	if err = a.DB.Ping(); err != nil {
		log.Fatalln("Connecting to DB failed: ", err)
	}
	//defer a.DB.Close() // This seems to close DB as soon as App.init() is complete.

	a.Router = mux.NewRouter()
	a.initRoutes()
}

// initRoutes() creates all the required API routes
func (a *App) initRoutes() {
	a.Router.HandleFunc("/v1/device/{orgID}/waitingroom", a.listWait).Methods("GET")
	a.Router.HandleFunc("/v1/device/waitingroom", a.addWait).Methods("POST")
	a.Router.HandleFunc("/v1/device/profile", a.getProfile).Methods("GET")
}

// run() starts the API server
func (a *App) run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
