package utils

import (
	"html/template"
	"net/http"
)

// GeneralPage is struct to define normal page
type GeneralPage struct {
	PageData   interface{}
	CurrentURL string
}

// GetTemplate return template var
func GetTemplate() *template.Template {
	tpl := template.Must(template.ParseGlob("templates/*.html"))
	return tpl
}

// GetPageStructure return populated general page structure
func GetPageStructure(w http.ResponseWriter, r *http.Request) *GeneralPage {
	page := GeneralPage{PageData: "", CurrentURL: r.URL.RequestURI()}
	return &page
}
