package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type waitingList struct {
	DevID     string `json:"device_id,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
}

type profile struct {
	ProfileURL string `json:"profile_url,omitempty"`
}

// App.listWait() is used to list an org's devices that are currenty in the waiting room.
func (a *App) listWait(w http.ResponseWriter, r *http.Request) {
	query := `SELECT device_id, created_at FROM devices`
	rows, err := a.DB.Query(query)
	var res []waitingList
	if err != nil {
		log.Printf("listWait(): Error getting waiting room devices: %v\n", err)
		http.Error(w, string("error getting devices"), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var item waitingList
		err = rows.Scan(&item.DevID, &item.CreatedAt)
		if err != nil {
			log.Printf("listWait(): Error getting next item: %v\n", err)
			http.Error(w, string("error getting devices"), http.StatusInternalServerError)
			return
		}
		res = append(res, item)
	}

	writeJSONResponse(w, http.StatusOK, res)
}

// App.addWait() is used by a newly onboarded device to periodically check for an assigned Profile URL
func (a *App) addWait(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	enableCors(&w)

	// Retrieve context values from authorizer middleware
	pathVars, ok := r.Context().Value(a.PathVarsKey).(map[string]string)
	if !ok {
		log.Println("Error asserting pathVarsKey context to string.")
		http.Error(w, "getting endpoint context", http.StatusInternalServerError)
		return
	}
	orgID := pathVars["orgID"]

	// get deviceID from context
	devID, ok := r.Context().Value(a.DevKey).(string)
	if !ok {
		log.Println("Error asserting devID context to string.")
		log.Println(r.Context().Value(a.DevKey))
		http.Error(w, "error getting endpoint context", http.StatusInternalServerError)
		return
	}

	// Run update query
	query := `
		UPDATE devices
		SET joined_at = ?
		WHERE device_id = ?
		AND org_id = ?;`
	res, err := a.DB.Exec(query, time.Now().Unix(), devID, orgID)
	if err != nil {
		log.Printf("App.addWait(): Error running update query: %v\n", err)
		http.Error(w, "Adding device to waitlist", http.StatusInternalServerError)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("App.addWait(): Error returning number rows affected: %v\n", err)
	}
	if rows != 1 {
		log.Printf("App.addWait(): Device not updated properly. Rows affected: %d, should be 1.\n", rows)
		http.Error(w, "Adding device to waitlist", http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, nil)
}

// App.getProfile() is used by devices to periodically check for a PRofile URL until one is provided.
func (a *App) getProfile(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	enableCors(&w)

	// Retrieve context values from authorizer middleware
	pathVars, ok := r.Context().Value(a.PathVarsKey).(map[string]string)
	if !ok {
		log.Println("Error asserting pathVarsKey context to string.")
		http.Error(w, "getting endpoint context", http.StatusInternalServerError)
		return
	}
	orgID := pathVars["orgID"]

	// get deviceID from context
	deviceID, ok := r.Context().Value(a.DevKey).(string)
	if !ok {
		log.Println("Error asserting devID context to string.")
		log.Println(r.Context().Value(a.DevKey))
		http.Error(w, "error getting endpoint context", http.StatusInternalServerError)
		return
	}

	// check for profile and return it
	var profURL profile
	query := `
		SELECT profile_url
		FROM devices
		WHERE device_id = ?
		AND org_id = ?`
	rows, err := a.DB.Query(query, deviceID, orgID)
	if err != nil {
		log.Printf("getKey(): Error querying for device_key: %v\n", err)
		writeJSONResponse(w, http.StatusOK, nil)
	}
	ok = rows.Next()
	if !ok {
		log.Printf("App.getProfile(): Error getting next row: %v", err)
		writeJSONResponse(w, http.StatusOK, nil)
	}
	err = rows.Scan(&profURL.ProfileURL)
	if err != nil {
		log.Printf("App.getProfile(): error retrieving row: %v\n", err)
		writeJSONResponse(w, http.StatusOK, nil)
	}

	// update last_seen timestamp
	query = `
		UPDATE devices
		SET last_seen = ?
		WHERE device_id = ?
		AND org_id = ?;`
	res, err := a.DB.Exec(query, time.Now().Unix(), deviceID, orgID)
	if err != nil {
		log.Printf("App.addWait(): Error running update query: %v\n", err)
		http.Error(w, "Adding device to waitlist", http.StatusInternalServerError)
		return
	}
	retRows, err := res.RowsAffected()
	if err != nil {
		log.Printf("App.getProfile(): Error returning number rows affected: %v\n", err)
	}
	if retRows != 1 {
		log.Printf("App.getProfile(): Device not updated properly. Rows affected: %d, should be 1.\n", retRows)
		http.Error(w, "Adding device to waitlist", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, profURL) // TODO: need to return profile URL object
}
