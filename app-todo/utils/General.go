package utils

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	// dbo object
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

// GeneralPage is struct to define normal page
type GeneralPage struct {
	PageData     interface{}
	FlashMessage interface{}
	CurrentURL   string
	Pager        interface{}
}

// FLASH_SESSION to store flash messages
const FLASH_SESSION string = "flash-message"

// StatusComplete status complete for todo
const StatusComplete = 1

// StatusIncomplete status complete for todo
const StatusIncomplete = 0

// GetDBO return db object
func GetDBO() *sql.DB {
	db, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

// GetTemplate return template object
func GetTemplate() *template.Template {
	tpl := template.Must(template.ParseGlob("templates/*.html"))
	return tpl
}

// GetPageStructure return populated general page structure
func GetPageStructure(w http.ResponseWriter, r *http.Request) (*GeneralPage, *sessions.Session) {
	_, flashSession := GetCookieStore(r, FLASH_SESSION)
	fm := flashSession.Flashes("message")
	flashSession.Save(r, w)
	page := GeneralPage{PageData: "", FlashMessage: fm, CurrentURL: r.URL.RequestURI()}
	return &page, flashSession
}

// GetCookieStore return cookie store
func GetCookieStore(r *http.Request, sessionName string) (*sessions.CookieStore, *sessions.Session) {
	store := sessions.NewCookieStore([]byte("a-secret-string"))
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Fatalln(err)
	}
	return store, session
}
