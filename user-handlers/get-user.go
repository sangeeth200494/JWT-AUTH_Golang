package userhandlers

import "net/http"

func RetrieveUserByID(w http.ResponseWriter, r *http.Request) {
	// It informs the client (browser, Postman, frontend app, etc.) that the response body will be in JSON format.
	w.Header().Set("Content-Type", "application/json")

}
