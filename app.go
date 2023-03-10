package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type key int

// App has router and db instances
type App struct {
	DevRouter   *mux.Router
	DB          *sql.DB
	DevKey      key
	PathVarsKey key
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

	// Set context id
	a.DevKey = 1
	a.PathVarsKey = 2

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

	a.DevRouter = mux.NewRouter().PathPrefix("/v1/device").Subrouter()
	a.initRoutes()
}

// authShim() is a middleware substitute for actual authorizer code
func (a *App) authShim(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		enableCors(&w)

		// authenticate 'Authorization' token
		var deviceID string

		// Check for missing Authorization header
		if r.Header["Authorization"] == nil {
			http.Error(w, string("missing Authorization header"), http.StatusUnauthorized)
			log.Println("missing Authorization header")
			return
		}
		tokenStr := r.Header["Authorization"][0]

		// Parse and validate JWT from Authorization header
		token, err := jwt.Parse(tokenStr, getKey, jwt.WithValidMethods([]string{"HS512"}))
		if err != nil {
			http.Error(w, string("invalid Authorization"), http.StatusForbidden)
			log.Printf("JWT parse returned error: %v", err)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		deviceID = fmt.Sprintf("%v", claims["device_id"]) // TODO: lookup key using this deviceID
		log.Println("DeviceID: ", deviceID)               // TODO: remove this message
		// Create context in order to pass back device key and path variables
		ctx := context.WithValue(context.TODO(), a.DevKey, deviceID)
		ctx = context.WithValue(ctx, a.PathVarsKey, mux.Vars(r))
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

// initRoutes() creates all the required API routes
func (a *App) initRoutes() {
	// Device endpoints
	a.DevRouter.Methods("GET").Path("/{orgID}/waiting_room").HandlerFunc(a.listWait)
	a.DevRouter.Methods("POST").Path("/{orgID}/waiting_room").HandlerFunc(a.authShim(a.addWait))
	a.DevRouter.Methods("GET").Path("/{orgID}/profile").HandlerFunc(a.authShim(a.getProfile))

	// Webapp endpoints
}

// run() starts the API server
func (a *App) run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.DevRouter))
}
