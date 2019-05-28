package user

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/deepaksinghkushwah/app-blog/utils"
)

var tpl *template.Template

func init() {
	tpl = utils.GetTemplate()
}

// LoginGetHandler show form
func LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := utils.GetPageStructure(w, r)
	if page.IsLoggedIn {
		http.Redirect(w, r, "/secret", http.StatusSeeOther)
	}

	err := tpl.ExecuteTemplate(w, "login.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// LoginPostHandler process form
func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)
	_, userSession := utils.GetCookieStore(r, utils.STORE_SESSION)

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	db := utils.GetDB()
	defer db.Close()
	var dbUsername, dbEmail, dbFullname string
	var dbUserID int
	err := db.QueryRow("SELECT id, username, email, fullname FROM `user` WHERE username = ? AND password = ?", username, password).Scan(&dbUserID, &dbUsername, &dbEmail, &dbFullname)
	if err != nil {
		if err == sql.ErrNoRows {
			flashSession.AddFlash("Invalid username or password", "message")
			userSession.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
	if dbUsername != "" {
		//fmt.Println("Here at loggedin")
		userSession.Values["isLoggedIn"] = true
		userSession.Values["userID"] = dbUserID
		userSession.Save(r, w)
		http.Redirect(w, r, "/secret", 302)
	}

}

// RegisterGetHandler show register form
func RegisterGetHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := utils.GetPageStructure(w, r)
	if page.IsLoggedIn {
		http.Redirect(w, r, "/secret", http.StatusSeeOther)
	}

	err := tpl.ExecuteTemplate(w, "register.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// RegisterPostHandler process register form
func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	defer db.Close()
	_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)

	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	fullname := r.PostForm.Get("fullname")
	email := r.PostForm.Get("email")
	if username == "" || password == "" || fullname == "" || email == "" {
		flashSession.AddFlash("All fields are required", "message")
		flashSession.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}

	var existingUsername string
	err := db.QueryRow("SELECT username FROM `user` WHERE username='" + username + "' OR email = '" + email + "'").Scan(&existingUsername)

	if err != nil {
		if err == sql.ErrNoRows {
			_, err := db.Exec("INSERT INTO `user` (username, password, email, fullname) VALUES (?,?,?,?)", username, password, email, fullname)
			if err != nil {
				log.Fatal(err)
			} else {
				flashSession.AddFlash("Username "+existingUsername+" registered successfully", "message")
				flashSession.Save(r, w)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		} else {
			log.Fatalln(err)
		}
	}
	if existingUsername != "" {
		flashSession.AddFlash("Username "+existingUsername+" already exists,please choose another", "message")
		flashSession.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}
	//http.Redirect(w, r, "/register", http.StatusSeeOther)
}

// LogoutHandler logout user from session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_, userSession := utils.GetCookieStore(r, utils.STORE_SESSION)
	userSession.Values["isLoggedIn"] = false
	userSession.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
