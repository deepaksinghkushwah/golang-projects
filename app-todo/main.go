package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/deepaksinghkushwah/app-todo/pagination"
	"github.com/deepaksinghkushwah/app-todo/utils"
	"github.com/gorilla/mux"
)

var tpl *template.Template
var perPage = 10

//Todo struct
type Todo struct {
	ID        int
	Title     string
	CreatedAt string
	UpdatedAt string
	Status    string
}

func init() {
	tpl = utils.GetTemplate()
}
func main() {
	r := setupRouter()
	// following line should be called after setupRouter function, otherwise system not work
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":5000", r)
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/", HomePostHandler).Methods("POST")
	r.HandleFunc("/change-status", ChangeStatusHandler).Methods("GET")
	return r
}

//HomeHandler home function
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	/*db := utils.GetDBO()
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec("INSERT INTO `todo`(title, created_at, status) VALUES(?,?,?)", "Sample Todo123", createdAt, utils.StatusIncomplete)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()*/
	page, _ := utils.GetPageStructure(w, r)
	var todos []Todo
	db := utils.GetDBO()
	defer db.Close()

	q := "SELECT count(id) FROM todo ORDER BY id DESC"
	var totalRows int
	err := db.QueryRow(q).Scan(&totalRows)
	//fmt.Println("Total rows : ", totalRows)
	if err != nil {
		log.Fatalln(err)
	}
	var offset int

	currentPage := 0
	if r.URL.Query().Get("page") != "" {
		currentPage, _ = strconv.Atoi(r.URL.Query().Get("page"))
		offset = (currentPage - 1) * perPage
	} else {
		currentPage = 0
		offset = 0
	}

	url := "/blog/list"
	pager := pagination.New(totalRows, perPage, currentPage, url)
	page.Pager = pager

	rows, err := db.Query("SELECT id, title, created_at, updated_at, status FROM todo ORDER BY id DESC limit ?,?", offset, perPage)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Error: no results")
		} else {
			log.Fatalln(err)
		}
	} else {
		for rows.Next() {
			var id, status int
			var title, istatus string
			var createdAt, updatedAt time.Time
			rows.Scan(&id, &title, &createdAt, &updatedAt, &status)
			if status == 1 {
				istatus = "Completed"
			} else {
				istatus = "Incomplete"
			}
			todos = append(todos, Todo{ID: id, Title: title, CreatedAt: createdAt.Format(time.RFC1123), UpdatedAt: updatedAt.Format(time.RFC1123), Status: istatus})
		}
	}

	page.PageData = todos
	err = tpl.ExecuteTemplate(w, "index.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

//HomePostHandler add todo
func HomePostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	if title == "" {
		_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)
		flashSession.AddFlash("Error at adding todo, empty title string!!!", "message")
		flashSession.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		createdAt := time.Now().Format("2006-01-02 15:04:05")
		status := 0
		db := utils.GetDBO()
		defer db.Close()
		_, err := db.Exec("INSERT INTO todo (title, created_at, status) VALUES(?,?,?)", title, createdAt, status)
		if err != nil {
			log.Fatalln(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func ChangeStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var dbStatus bool
	db := utils.GetDBO()
	defer db.Close()
	err := db.QueryRow("SELECT `status` FROM todo WHERE id = ?", id).Scan(&dbStatus)

	if err != nil {
		if err == sql.ErrNoRows {

		} else {
			log.Fatalln(err)
		}
	} else {
		newStatus := !dbStatus
		updatedAt := time.Now().Format("2006-01-02 15:04:05")
		_, err := db.Exec("UPDATE todo SET updated_at = ?, status = ? WHERE id = ?", updatedAt, newStatus, id)
		if err != nil {
			log.Fatalln(err)
		}
		_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)
		flashSession.AddFlash("Status Updated", "message")
		flashSession.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
