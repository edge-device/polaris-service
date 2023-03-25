package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSONResponse() is helper that returns JSON HTTP response
func writeJSONResponse(w http.ResponseWriter, respCode int, payload interface{}) {
	p, err := json.Marshal(payload)
	if err != nil {
		log.Printf("writeJSONResponse(): JSON marshal failed. %v\n", err)
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(respCode)
	w.Write(p)
}

// genRand generates and returns hash of random []byte
// func genRand(c int) (string, error) {
// 	b := make([]byte, c)
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		log.Println("error generating random []byte:", err)
// 		return "", err
// 	}
// 	encoded := base64.StdEncoding.EncodeToString(b)
// 	return encoded, nil
// }

// Enable CORS support on server
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
