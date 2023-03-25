package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

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
	Conf        config
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

	// Create DB health check loop
	var err error
	dbFailure := true
	for i := 1; i < 10; i++ {
		a.DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("DB fail to open loop#%d: ", i)
			duration := time.Second * 5
			time.Sleep(duration)
			continue
		}
		if err = a.DB.Ping(); err != nil {
			log.Printf("DB fail to connect loop#%d: ", i)
			duration := time.Second * 5
			time.Sleep(duration)
			continue
		}
		dbFailure = false               // If execution gets here, DB is good
		log.Println("DB must be good!") // TODO: remove this debug message
		break
	}
	if dbFailure {
		log.Fatalln("Connecting to DB failed: ", err)
	}

	a.DevRouter = mux.NewRouter().PathPrefix("/v1/device").Subrouter()
	a.initRoutes()
}

// authDevice() is a middleware substitute for actual authorizer code
func (a *App) authDevice(h http.HandlerFunc) http.HandlerFunc {
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
		token, err := jwt.Parse(tokenStr, a.getKey, jwt.WithValidMethods([]string{"HS512"}))
		if err != nil {
			http.Error(w, string("invalid Authorization"), http.StatusForbidden)
			log.Printf("JWT parse returned error: %v", err)
			return
		}

		// Get device_id from JWT claims
		claims := token.Claims.(jwt.MapClaims)
		deviceID = fmt.Sprintf("%v", claims["device_id"])

		// Create context in order to pass back device key and path variables
		ctx := context.WithValue(context.TODO(), a.DevKey, deviceID)
		ctx = context.WithValue(ctx, a.PathVarsKey, mux.Vars(r))
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

// App.getKey() is used by device authenticators to lookup a device's key
func (a *App) getKey(token *jwt.Token) (interface{}, error) {
	// Extract deviceID & orgID claims from token
	claims := token.Claims.(jwt.MapClaims)
	deviceID := fmt.Sprintf("%v", claims["device_id"])
	orgID := fmt.Sprintf("%v", claims["org_id"])
	log.Println("DeviceID(getKey): ", deviceID) // TODO: remove this message

	// Get device key from DB
	var deviceKey string
	query := `
		SELECT device_key
		FROM devices
		WHERE device_id = ?
		AND org_id = ?`
	rows, err := a.DB.Query(query, deviceID, orgID)
	if err != nil {
		log.Printf("getKey(): Error querying for device_key: %v\n", err)
		return nil, fmt.Errorf("getKey(): Error querying for device_key: %w", err)
	}
	ok := rows.Next()
	if !ok {
		log.Printf("getKey(): Error getting next row: %v", err)
		return nil, fmt.Errorf("getKey(): Error getting next row: %w", err)
	}
	err = rows.Scan(&deviceKey)
	if err != nil {
		log.Printf("getKey(): error retrieving row: %v\n", err)
		return nil, fmt.Errorf("getKey(): error retrieving row: %w", err)
	}

	devKeyBin, err := base64.StdEncoding.DecodeString(deviceKey)
	if err != nil {
		log.Printf("getKey(): error decoding deviceKey: %v\n", err)
		return nil, fmt.Errorf("getKey(): error decoding deviceKey: %w", err)
	}
	log.Println("Retrieved deviceKey:", deviceKey) // TODO: remove this debug message

	return devKeyBin, nil
}

// App.initRoutes() creates all the required API routes
func (a *App) initRoutes() {
	// Device endpoints
	a.DevRouter.Methods("GET").Path("/{orgID}/waiting_room").HandlerFunc(a.listWait)
	a.DevRouter.Methods("POST").Path("/{orgID}/waiting_room").HandlerFunc(a.authDevice(a.addWait))
	a.DevRouter.Methods("GET").Path("/{orgID}/profile").HandlerFunc(a.authDevice(a.getProfile))
	a.DevRouter.Methods("GET").Path("/auth/token").HandlerFunc(a.getTokens)
	a.DevRouter.Methods("POST").Path("/auth/token/refresh").HandlerFunc(a.refreshTokens)
	a.DevRouter.Methods("GET").Path("/auth/join").HandlerFunc(a.createAccount)

	// Webapp endpoints
}

// run() starts the API server
func (a *App) run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.DevRouter))
}
