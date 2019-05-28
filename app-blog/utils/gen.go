package utils

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	// added for mysql support
	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/sessions"
)

// FLASH_SESSION to store flash messages
const FLASH_SESSION string = "flash-message"

// STORE_SESSION to store site vars
const STORE_SESSION string = "store"

// GeneralPage is struct to define normal page
type GeneralPage struct {
	IsLoggedIn   bool
	PageData     interface{}
	FlashMessage interface{}
	CurrentURL   string
	UserID       int
	Pager        interface{}
}

// GetTemplate return template var
func GetTemplate() *template.Template {
	tpl := template.Must(template.ParseGlob("templates/*.html"))
	return tpl
}

// GetDB return db object
func GetDB() *sql.DB {
	db, err := sql.Open("mysql", "root:deepak@tcp(127.0.0.1:3306)/test2?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	//db.SetConnMaxLifetime(time.Hour)
	//db.SetMaxOpenConns(2)
	//db.SetMaxIdleConns(1)
	return db
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

// CheckLoggedIn to check if user session is set
func CheckLoggedIn(r *http.Request) bool {
	_, session := GetCookieStore(r, STORE_SESSION)
	var x interface{} = session.Values["isLoggedIn"]
	if x == nil {
		x = false
	}
	return x.(bool)
}

// GetLogginUserID to return login id or loggedin user or return 0
func GetLogginUserID(r *http.Request) int {
	_, session := GetCookieStore(r, STORE_SESSION)
	x := session.Values["userID"]
	if x == nil {
		x = 0
	}
	u, _ := x.(int)
	return u
}

// GetPageStructure return populated general page structure
func GetPageStructure(w http.ResponseWriter, r *http.Request) (*GeneralPage, *sessions.Session) {
	isLoggedIn := CheckLoggedIn(r)
	_, flashSession := GetCookieStore(r, FLASH_SESSION)
	fm := flashSession.Flashes("message")
	flashSession.Save(r, w)
	page := GeneralPage{IsLoggedIn: isLoggedIn, PageData: "", FlashMessage: fm, CurrentURL: r.URL.RequestURI(), UserID: GetLogginUserID(r)}
	return &page, flashSession
}

// LoginRequired redirect user if user is not logged in
func LoginRequired(w http.ResponseWriter, r *http.Request) {
	page, flashSession := GetPageStructure(w, r)
	if page.IsLoggedIn == false {
		flashSession.AddFlash("You are not authorized to perform this action!!!", "message")
		flashSession.Save(r, w)
		http.Redirect(w, r, "/login", 301)
	}
}

func RedirectWithMessage(w http.ResponseWriter, r *http.Request, msg string) {
	_, flashSession := GetCookieStore(r, FLASH_SESSION)
	flashSession.AddFlash(msg, "message")
	flashSession.Save(r, w)
	http.Redirect(w, r, "/blog/list", http.StatusSeeOther)
}
