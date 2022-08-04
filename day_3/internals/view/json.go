package view

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func RenderJSON(w http.ResponseWriter, v interface{}, template *template.Template) {
	js, err := json.Marshal(v)
	if err != nil {
		HandleErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func HandleErrorJSON(w http.ResponseWriter, msg string, code int) {
	type errorJSON struct {
		Error string `json:"error"`
	}
	json, _ := json.Marshal(errorJSON{msg})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(json)
}
