package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
)

type waitingList struct {
	DevID     string `json:"device_id,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
}

type profile struct {
	ProfileURL string `json:"profile_url,omitempty"`
}

type tknObject struct {
	AccTkn  string `json:"access_token"`
	TknType string `json:"token_type"`
	Scope   string `json:"scope"`
}

type emailObject struct {
	Email      string `json:"email"`
	Verified   bool   `json:"verified"`
	Primary    bool   `json:"primary"`
	Visibility string `json:"visibility"`
}

type tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type refresh struct {
	UID          string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
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

// App.getTokens() is used to request access & refresh tokens using callback code from
// Github Oauth requests.
func (a *App) getTokens(w http.ResponseWriter, r *http.Request) {
	accessCode := r.URL.Query().Get("code")
	clientID := "d46f377df25a400e9c03"
	clientSecret := "9c0d7c235afbd8847b82e8bbe6a8d2fc9aea335d"

	// Exchange code from callback for access token
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr}
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, accessCode)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		log.Printf("App.getTokens(): could not create Github access token request: %v\n", err)
		writeJSONResponse(w, http.StatusUnauthorized, nil)
		return
	}
	req.Header.Set("accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Printf("App.getTokens(): could not send Github access token request: %v\n", err)
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("App.getTokens(): could not retrieve body for token request: %v\n", err)
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}

	// Using access token, request user's email addresses
	tkn := tknObject{}
	err = json.Unmarshal(b, &tkn)
	if err != nil {
		log.Printf("App.getTokens(): could not unmarshall token request: %v\n", err)
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}
	reqURL = "https://api.github.com/user/emails"
	req, err = http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		log.Printf("App.getTokens(): could not create Github email request: %v\n", err)
		writeJSONResponse(w, http.StatusUnauthorized, nil)
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tkn.AccTkn))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	res, err = client.Do(req)
	if err != nil {
		log.Printf("App.getTokens(): could not send Github get email request: %v\n", err)
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		log.Printf("App.getTokens(): could not retrieve body for email request: %v\n", err)
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}
	emailAddrs := []emailObject{}
	err = json.Unmarshal(b, &emailAddrs)
	if err != nil {
		log.Printf("App.getTokens(): could not unmarshall email address: %v\n", err)
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}

	// Get primary email address
	var emailAddr string
	for _, item := range emailAddrs {
		if item.Primary {
			emailAddr = strings.ToLower(item.Email)
			break
		}
	}

	// Verify user and request tokens
	err = a.verifyUser(emailAddr)
	if err != nil {
		log.Printf("App.getTokens(): verify email error: %v\n", err)
		writeJSONResponse(w, http.StatusUnauthorized, nil)
		return
	}
	aTkn, rTkn, err := a.createTokens(emailAddr)
	if err != nil {
		log.Printf("App.getTokens(): issue generating tokens: %v\n", err)
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}

	payload := tokens{aTkn, rTkn}
	writeJSONResponse(w, http.StatusOK, payload)
}

// App.refreshTokens() creates new access and refresh token
func (a *App) refreshTokens(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	enableCors(&w)

	// Retrieve http request body
	var b refresh
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("App.refreshTokens(): could not retrieve body for token refresh: ", err)
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}
	// Parse http request JSON body
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Println("App.refreshTokens(): could not unmarshal body: ", err)
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}

	// Parse current refresh token
	refreshToken, err := jwt.Parse(b.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("JWT: Signing method expecting HMAC': %v", token.Header["alg"])
		}
		return a.Conf.signKey, nil
	})
	if err != nil {
		log.Println("App.refreshTokens(): refresh token validation failed: ", err)
		writeJSONResponse(w, http.StatusUnauthorized, nil)
		return
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || !refreshToken.Valid {
		log.Println("App.refreshTokens(): refresh token claims validation failed: ")
		writeJSONResponse(w, http.StatusUnauthorized, nil)
		return
	}
	if err := claims.Valid(); err != nil {
		log.Println("App.refreshTokens(): refresh token expired: ", err)
		writeJSONResponse(w, http.StatusUnauthorized, nil)
		return
	}
	if claims["type"] != "refresh" {
		log.Println("App.refreshTokens(): token not type 'refresh' ", err)
		writeJSONResponse(w, http.StatusUnauthorized, nil)
		return
	}
	currentTokenID := strconv.FormatFloat(claims["tid"].(float64), 'f', 0, 64)

	// Ensure refresh token not revoked
	var expire int64
	query := `
		SELECT expires_at 
		FROM access_token 
		WHERE token_id = ?;`
	err = a.DB.QueryRow(query, currentTokenID).Scan(&expire)
	switch {
	case err == sql.ErrNoRows:
		log.Println("App.refreshTokens(): refresh token ID not found, maybe revoked", err)
		writeJSONResponse(w, http.StatusUnauthorized, nil)
		return
	case err != nil:
		log.Println("App.refreshTokens(): TokenID query failed", err)
		writeJSONResponse(w, http.StatusUnauthorized, nil)
		return
	}
	if expire < time.Now().Unix() {
		log.Printf("App.refreshTokens(): token %s expired", currentTokenID)
		writeJSONResponse(w, http.StatusUnauthorized, nil)
		return
	}

	// Create new tokens
	aTkn, rTkn, err := a.createTokens(claims["user"].(string))
	if err != nil {
		log.Printf("App.getTokens(): issue generating tokens: %v\n", err)
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}

	payload := tokens{aTkn, rTkn}
	writeJSONResponse(w, http.StatusOK, payload)
}

// App.verifyUser() updates user last login. Returns error if user not found.
func (a *App) verifyUser(uid string) error {
	query := `
		UPDATE users
		SET last_login = ?
		WHERE user_id = ?`
	res, err := a.DB.Exec(query, time.Now().Unix(), uid)
	if err != nil {
		return fmt.Errorf("App.verifyUser(): updating user. %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows < 1 {
		return fmt.Errorf("App.verifyUser(): user not found")
	}

	return nil
}

// App.addUser() creates user account if not already exist
// func (a *App) addUser(uid string) error {
// 	// First check if user already exists
// 	query := `
// 		SELECT * FROM users
// 		WHERE user_id = ?`
// 	res, err := a.DB.Exec(query, uid)
// 	if err != nil {
// 		return fmt.Errorf("App.addUser(): checking if user exists. %w", err)
// 	}
// 	rows, _ := res.RowsAffected()
// 	if rows == 1 {
// 		log.Println("App.addUser(): user exists already")
// 		return nil
// 	}

// 	// TODO: add user here

// 	return nil
// }

// App.createTokens() is used to create and sign access and refresh tokens
func (a *App) createTokens(uid string) (accessTkn, refreshTkn string, err error) {
	// TODO: These duration times need to be configured by configuration, not hard coded
	timeNow := time.Now()
	accessExp := timeNow.Add(time.Second * 10)
	refreshExp := timeNow.Add(time.Hour * 24 * 10)
	uid = strings.ToLower(uid)

	// Create token DB entry and retrieve auto-generated token ID
	query := `
		INSERT INTO access_token (user_id, created_at, expires_at)
		VALUES(?, ?, ?);`
	res, err := a.DB.Exec(query, uid, timeNow.Unix(), refreshExp.Unix())
	if err != nil {
		err = fmt.Errorf("App.createTokens() running insert token query. %w", err)
		return accessTkn, refreshTkn, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		err = fmt.Errorf("App.createTokens() running insert token query rows affected. %w", err)
		return accessTkn, refreshTkn, err
	}
	if rows != 1 {
		err = fmt.Errorf("App.createTokens() Rows affected: %d, should be 1. %w", rows, err)
		return accessTkn, refreshTkn, err
	}
	tokenID, err := res.LastInsertId()
	if err != nil {
		err = fmt.Errorf("App.createTokens() getting token ID. %w", err)
		return accessTkn, refreshTkn, err
	}

	// TODO:  Delete current refresh token

	// Delete expired refresh tokens
	t := time.Now().Unix()
	_, err = a.DB.Exec("DELETE FROM access_token WHERE expires_at < ?", t)
	if err != nil {
		log.Println("issue deleting expired tokens", err)
	}

	// Generate Access and Refresh tokens
	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  timeNow.Unix(),
		"exp":  accessExp.Unix(),
		"user": uid,
		"type": "access",
		"tid":  tokenID,
	})
	rtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  timeNow.Unix(),
		"exp":  refreshExp.Unix(),
		"user": uid,
		"type": "refresh",
		"tid":  tokenID,
	})

	// Sign tokens with private key
	accessTkn, err = atoken.SignedString(a.Conf.signKey)
	if err != nil {
		err = fmt.Errorf("App.createTokens(): signing access token. %w", err)
		return accessTkn, refreshTkn, err
	}
	refreshTkn, err = rtoken.SignedString(a.Conf.signKey)
	if err != nil {
		err = fmt.Errorf("App.createTokens(): signing refresh token. %w", err)
		return accessTkn, refreshTkn, err
	}

	return accessTkn, refreshTkn, err
}

func (a *App) createAccount(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	query := `
		INSERT INTO users (user_id, created_at)
		VALUES(?, ?);`
	res, err := a.DB.Exec(query, email, time.Now().Unix())
	if err != nil {
		log.Println("App.createAccount() inserting new user.", err)
		writeJSONResponse(w, http.StatusBadRequest, nil)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("App.createAccount() problem getting rows affected", err)
		writeJSONResponse(w, http.StatusBadRequest, nil)
		return
	}
	if rows != 1 {
		log.Println("App.createAccount() number rows inserted not = 1.", err)
		writeJSONResponse(w, http.StatusBadRequest, nil)
		return
	}

	writeJSONResponse(w, http.StatusOK, nil)
}
