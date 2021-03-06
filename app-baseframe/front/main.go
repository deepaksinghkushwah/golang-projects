package main

import (
	"html/template"
	"log"
	"net/http"

	"baseframe/front/user"
	utils "baseframe/utils/general"

	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	var err error

	if err != nil {
		log.Fatalln(err)
	}
	tpl = utils.GetTemplate()
}

func main() {

	//runtime.GOMAXPROCS(0)
	r := setupRoutes()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":5000", r)
}

func setupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/login", user.LoginGetHandler).Methods("GET")
	r.HandleFunc("/login", user.LoginPostHandler).Methods("POST")

	r.HandleFunc("/register", user.RegisterGetHandler).Methods("GET")
	r.HandleFunc("/register", user.RegisterPostHandler).Methods("POST")

	r.HandleFunc("/secret", secretHandler).Methods("GET")
	r.HandleFunc("/logout", user.LogoutHandler).Methods("GET")

	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := utils.GetPageStructure(w, r)
	err := tpl.ExecuteTemplate(w, "index.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

func secretHandler(w http.ResponseWriter, r *http.Request) {
	page, flashSession := utils.GetPageStructure(w, r)
	if page.IsLoggedIn {
		err := tpl.ExecuteTemplate(w, "secret.html", page)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		flashSession.AddFlash("You are not authorized to perform this action!!!", "message")
		flashSession.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
