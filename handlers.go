package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// type returnStatus struct {
// 	Error   bool   `json:"error,omitempty"`
// 	Code    int    `json:"code,omitempty"`
// 	Type    string `json:"type,omitempty"`
// 	Message string `json:"message,omitempty"`
// }

// type paging struct {
// 	Cursor string `json:"next_cursor"`
// 	Qty    int    `json:"quantity"`
// }

// type returnMsg struct {
// 	Status returnStatus `json:"status,omitempty"`
// 	Paging *paging      `json:"paging,omitempty"`
// 	Data   interface{}  `json:"data,omitempty"`
// }

type waitingList struct {
	DevID     string `json:"device_id,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
}

func (a *App) listWait(w http.ResponseWriter, r *http.Request) {
	query := `SELECT device_id, created_at FROM waiting_room`
	rows, err := a.DB.Query(query)
	if err != nil {
		log.Printf("listWait(): Error getting waiting room devices: %v\n", err)
		http.Error(w, string("error getting devices"), http.StatusInternalServerError)
		return
	}

	var res []waitingList
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

// App.addWait() is used by a new device to join the waiting room and check for Profile URL. Devices
// should periodically hit this endpoint to keep checking for a PRofile URL until one is provided.
func (a *App) addWait(w http.ResponseWriter, r *http.Request) {
	// TODO: Need to to first check if
	// query := `
	// 	INSERT INTO waiting_room (device_id, org_id, created_at, last_seen)
	// 	VALUES(?, ?, ?, ?);`
	// timeNow := time.Now()
	// res, err := a.DB.Exec(query, userID, timeNow.Unix(), timeNow.Unix())
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

// App.getProfile() is used by devices to periodically check for a PRofile URL until one is provided.
func (a *App) getProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Need to to first check if
	// query := `
	// 	INSERT INTO waiting_room (device_id, created_at, expires_at)
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

// writeJSONResponse() is helper that returns JSON HTTP response
func writeJSONResponse(w http.ResponseWriter, resCode int, payload interface{}) {
	p, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(resCode)
	w.Write(p)
}
