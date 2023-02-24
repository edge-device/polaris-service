package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	//"gorm.io/driver/mysql"
	//"gorm.io/gorm"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	//DB     *gorm.DB
}

// Initialize initializes the app with predefined configuration
func (a *App) init(config *config) {
	/* 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
	   		config.DB.Username,
	   		config.DB.Password,
	   		config.DB.Host,
	   		config.DB.Port,
	   		config.DB.Name,
	   		config.DB.Charset)

	   	db, err := gorm.Open(config.DB.Dialect, dsn)
	   	if err != nil {
	   		log.Fatal("Could not connect database")
	   	}

	   	a.DB = model.DBMigrate(db) */
	a.Router = mux.NewRouter()
	a.initRoutes()
}

// initRoutes() creates all the required API routes
func (a *App) initRoutes() {
	a.Router.HandleFunc("/devices", a.listDevices).Methods("GET")
}

// run() starts the API server
func (a *App) run(addr string) {
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
