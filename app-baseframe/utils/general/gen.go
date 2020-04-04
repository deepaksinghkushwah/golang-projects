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
	PageTitle    string
	IsLoggedIn   bool
	IsAdmin      bool
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
	db, err := sql.Open("mysql", "root:deepak@tcp(127.0.0.1:3306)/go_baseframe?parseTime=true")
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
	store := sessions.NewCookieStore([]byte("fsdf-@#$@wfsdf-sdfcvCXV$#V"))
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

// CheckAdminLoggedIn to check if admin user session is set
func CheckAdminLoggedIn(r *http.Request) bool {
	_, session := GetCookieStore(r, STORE_SESSION)
	var x interface{} = session.Values["isAdmin"]
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
	isAdmin := CheckAdminLoggedIn(r)
	_, flashSession := GetCookieStore(r, FLASH_SESSION)
	fm := flashSession.Flashes("message")
	flashSession.Save(r, w)
	page := GeneralPage{IsLoggedIn: isLoggedIn, IsAdmin: isAdmin, PageData: "", FlashMessage: fm, CurrentURL: r.URL.RequestURI(), UserID: GetLogginUserID(r)}
	return &page, flashSession
}

// LoginRequired redirect user if user is not logged in
func LoginRequired(w http.ResponseWriter, r *http.Request) {
	page, _ := GetPageStructure(w, r)
	if page.IsLoggedIn == false {
		RedirectWithMessage(w, r, "/login", "message", "You are not authorized to perform this action")
	}
}

// AdminLoginRequired redirect user if user is not logged in as admin
func AdminLoginRequired(w http.ResponseWriter, r *http.Request) {
	page, _ := GetPageStructure(w, r)
	//log.Println(page.IsAdmin, page.IsLoggedIn)
	if page.IsAdmin == false && page.IsLoggedIn == false {
		RedirectWithMessage(w, r, "/login", "message", "You are not authorized to perform this action")
	}
}

//RedirectWithMessage will redirect user with message on a page
func RedirectWithMessage(w http.ResponseWriter, r *http.Request, path string, msgCategory string, msg string) {
	_, flashSession := GetCookieStore(r, FLASH_SESSION)
	flashSession.AddFlash(msg, msgCategory)
	err := flashSession.Save(r, w)
	if err != nil {
		log.Fatalln(err.Error())
	}
	http.Redirect(w, r, path, http.StatusSeeOther)
}
