package lesson5

import (
	"encoding/json"
	"net/http"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "available",
		"version": app.Config.Version,
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		app.Logger.Println(err)
		http.Error(w, "A problem occured", http.StatusInternalServerError)
	}
}
