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
	// TODO: remove debug message below
	log.Printf("App.addWait(): orgID:%s\n", orgID)

	devID, ok := r.Context().Value(a.DevKey).(string)
	if !ok {
		log.Println("Error asserting devID context to string.")
		log.Println(r.Context().Value(a.DevKey))
		http.Error(w, "getting endpoint context", http.StatusInternalServerError)
		return
	}
	// TODO: remove debug message below
	log.Printf("App.addWait(): devID:%s\n", devID)

	// Run update query
	query := `
		UPDATE devices
		SET joined_at = ?
		WHERE device_id = ?
		AND org_id = ?;`
	log.Println(query) // TODO: remove debug message
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
	// TODO: Need to to first check if
	// query := `
	// 	INSERT INTO devices (device_id, created_at, expires_at)
	// 	VALUES(?, ?, ?);`
	// res, err := a.DB.Exec(query, userID, timeNow.Unix(), refreshExp.Unix())
	// if err != nil {
	// 	log.Printf("newToken(): Error running insert query: %v\n", err)
	// 	retStatus := returnStatus{true, 500, "Internal Server Error", "Internal Server Error"}
	// 	retMsg := apiRetMsg(retStatus, nil, nil)
	// 	http.Error(w, string(retMsg), http.StatusInternalServerError)
	// 	return
	// }
	// rows, err := res.RowsAffected()
	// if err != nil {
	// 	log.Printf("newToken(): Error returning number rows affected: %v\n", err)
	// }
	// if rows != 1 {
	// 	log.Printf("newToken(): Token ID not properly inserted into refresh table. Rows affected: %d, should be 1.", rows)
	// }
	writeJSONResponse(w, http.StatusOK, map[string]string{"result": "success"})
}
