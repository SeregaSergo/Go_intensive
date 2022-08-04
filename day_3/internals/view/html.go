package view

import (
	"bytes"
	"html/template"
	"net/http"
)

func RenderHtml(w http.ResponseWriter, v interface{}, template *template.Template) {
	var buf bytes.Buffer
	template.ExecuteTemplate(&buf, "page", v)
	w.Header().Set("Content-Type", "text/html")
	w.Write(buf.Bytes())
}
