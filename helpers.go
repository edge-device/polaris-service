package main

import (
	//"crypto/rand"
	//"encoding/base64"
	"encoding/json"
	//"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

// writeJSONResponse() is helper that returns JSON HTTP response
func writeJSONResponse(w http.ResponseWriter, resCode int, payload interface{}) {
	p, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(resCode)
	w.Write(p)
}

func getKey(token *jwt.Token) (interface{}, error) {
	return []byte("This is my secret key"), nil
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
