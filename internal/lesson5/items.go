package lesson5

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) FillFakeItems() {
	for i := 1; i <= 5; i++ {
		item := Item{
			ID:          int64(i),
			Name:        "Name" + strconv.Itoa(i),
			Description: "Description" + strconv.Itoa(i),
			Price:       int64(i),
			ImageLink:   "ImageLink" + strconv.Itoa(i),
		}
		app.Items = append(app.Items, item)
	}
}

func (app *Application) showItemsHandler(w http.ResponseWriter, r *http.Request) {
	if app.Items == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(app.Items)
	if err != nil {
		app.Logger.Println(err)
		http.Error(w, "A problem occured", http.StatusInternalServerError)
	}
}

func (app *Application) showSingleItemHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	if app.Items == nil || id > int64(len(app.Items)) {
		http.NotFound(w, r)
		return
	}

	var item Item
	found := false
	for i := 0; i < len(app.Items); i++ {
		if id == app.Items[i].ID {
			found = true
			item = app.Items[i]
			break
		}
	}

	if !found {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		app.Logger.Println(err)
		http.Error(w, "A problem occured", http.StatusInternalServerError)
	}
}
