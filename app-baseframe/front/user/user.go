package user

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	utils "github.com/deepaksinghkushwah/projects/app-baseframe/utils/general"
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
	//_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)
	_, userSession := utils.GetCookieStore(r, utils.STORE_SESSION)

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	db := utils.GetDB()
	defer db.Close()
	var dbUsername, dbEmail string
	var dbPassword []byte
	var dbUserID int
	err := db.QueryRow("SELECT id, username, password, email FROM `user` WHERE username = ?", username).Scan(&dbUserID, &dbUsername, &dbPassword, &dbEmail)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.RedirectWithMessage(w, r, "/login", "message", "Invalid username or password")
		} else {
			log.Fatalln(err)
		}
	}

	if dbUsername != "" {
		//log.Print()
		if err := bcrypt.CompareHashAndPassword(dbPassword, []byte(password)); err != nil {
			//utils.RedirectWithMessage(w, r, "Invalid username or password")
			log.Fatalln(err)
		} else {
			userSession.Values["isLoggedIn"] = true
			userSession.Values["userID"] = dbUserID
			userSession.Save(r, w)
			//http.Redirect(w, r, "/secret", 302)
			utils.RedirectWithMessage(w, r, "/secret", "message", "You are loggedin successfully")
		}

	} else {
		utils.RedirectWithMessage(w, r, "/login", "message", "Invalid username or password")
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
	//_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)

	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	email := r.PostForm.Get("email")
	if username == "" || password == "" || email == "" {
		utils.RedirectWithMessage(w, r, "/register", "message", "All fields are required")
	}

	var existingUsername string
	err := db.QueryRow("SELECT username FROM `user` WHERE username='" + username + "' OR email = '" + email + "'").Scan(&existingUsername)

	if err != nil {
		if err == sql.ErrNoRows {
			hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			_, err := db.Exec("INSERT INTO `user` (username, password, email, role_id) VALUES (?,?,?,?)", username, hash, email, 1)
			if err != nil {
				log.Fatal(err)
			} else {
				utils.RedirectWithMessage(w, r, "/register", "message", "Username "+existingUsername+" registered successfully")

			}
		} else {
			//log.Fatalln(err)
		}
	}
	if existingUsername != "" {
		utils.RedirectWithMessage(w, r, "/register", "message", "Username "+existingUsername+" already registered, please choose another")
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
