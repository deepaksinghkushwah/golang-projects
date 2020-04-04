package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/deepaksinghkushwah/projects/app-baseframe/admin/user"
	utils "github.com/deepaksinghkushwah/projects/app-baseframe/utils/general"
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

	http.ListenAndServe(":5001", r)
}

func setupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/login", user.LoginGetHandler).Methods("GET")
	r.HandleFunc("/login", user.LoginPostHandler).Methods("POST")
	r.HandleFunc("/dashboard", user.DashboardHandler).Methods("GET")
	r.HandleFunc("/logout", user.LogoutHandler).Methods("GET")

	// user add/edit/delete handlers
	//r.HandleFunc("/user/ct", user.Ct).Methods("GET")
	r.HandleFunc("/user/list", user.ListGetHandler).Methods("GET")
	r.HandleFunc("/user/create", user.CreateGetHandler).Methods("GET")
	r.HandleFunc("/user/create", user.CreatePostHandler).Methods("POST")
	r.HandleFunc("/user/edit/{id:[0-9]+}", user.EditGetHandler).Methods("GET")
	r.HandleFunc("/user/edit/{id:[0-9]+}", user.EditPostHandler).Methods("POST")
	r.HandleFunc("/user/delete/{id:[0-9]+}", user.DeleteGetHandler).Methods("GET")

	r.HandleFunc("/user/populate/{start:[0-9]+}/{end:[0-9]+}", user.PopulateUserSecond).Methods("GET")

	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := utils.GetPageStructure(w, r)
	err := tpl.ExecuteTemplate(w, "index.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}
