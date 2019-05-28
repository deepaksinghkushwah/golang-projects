package blog

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/deepaksinghkushwah/app-blog/pagination"
	"github.com/deepaksinghkushwah/app-blog/utils"
)

var tpl *template.Template
var perPage = 100

func init() {
	tpl = utils.GetTemplate()
}

// Blogs structure to define blog
type Blogs struct {
	ID        string
	Title     string
	Content   template.HTML
	CreatedAt string
	Author    int64
	Comments  []Comment
}

// Comment list of all comments for blog
type Comment struct {
	ID        int64
	BlogID    int64
	ParentID  int64
	Content   string
	CreatedAt string
	Author    Author
}

// Author user info who posted comment
type Author struct {
	ID   int64
	Name string
}

// List list all blogs
func List(w http.ResponseWriter, r *http.Request) {
	page, _ := utils.GetPageStructure(w, r)
	db := utils.GetDB()
	defer db.Close()
	/*if page.IsLoggedIn {
		http.Redirect(w, r, "/secret", http.StatusSeeOther)
	}*/
	var blogs []Blogs
	q := "SELECT count(id) FROM blog ORDER BY id DESC"
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

	result, err := db.Query("SELECT id, title, content, created_at, author FROM blog ORDER BY id DESC limit ?,?", offset, perPage)
	if err != nil {
		if err == sql.ErrNoRows {

		} else {
			log.Fatalln(err)
		}
	} else {
		for result.Next() {
			var id, title string
			var createdAt time.Time
			var content template.HTML
			var author int64
			result.Scan(&id, &title, &content, &createdAt, &author)
			d := createdAt.Format(time.RFC1123)
			//fmt.Println(createdAt, d)
			idInt, _ := strconv.Atoi(id)
			comments := GetComments(idInt)

			blogs = append(blogs, Blogs{ID: id, Title: title, Content: content, CreatedAt: d, Author: author, Comments: comments})
		}
	}

	page.PageData = blogs

	err = tpl.ExecuteTemplate(w, "blog-list.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetComments(blogID int) []Comment {
	var comments []Comment
	db := utils.GetDB()
	defer db.Close()
	query := "SELECT id, blog_id, parent_id, content, created_at, created_by FROM comment WHERE blog_id = ? ORDER BY id DESC"
	result, err := db.Query(query, blogID)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No comments found for blog", blogID)
		} else {
			log.Fatalln(err)
		}
	} else {
		for result.Next() {
			var id, parentID, dBlogID int64
			var createdAt time.Time
			var content string
			var createdBy int
			result.Scan(&id, &dBlogID, &parentID, &content, &createdAt, &createdBy)
			d := createdAt.Format(time.RFC1123)
			author := GetAuthor(createdBy)
			//fmt.Println(createdAt, d)

			comments = append(comments, Comment{ID: id, BlogID: dBlogID, ParentID: parentID, Content: content, CreatedAt: d, Author: author})
		}

	}
	return comments
}

// AddBlogGetHandler show blog form
func AddBlogGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LoginRequired(w, r)
	page, _ := utils.GetPageStructure(w, r)
	err := tpl.ExecuteTemplate(w, "blog-add.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// AddBlogPostHandler add blog post
func AddBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	utils.LoginRequired(w, r)
	r.ParseForm()
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	db := utils.GetDB()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO blog (title, content, created_at, status,author) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatalln(err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(title, content, time.Now().Format("2006-01-02 15:04:05"), 1, utils.GetLogginUserID(r))
	if err != nil {
		log.Fatalln(err)
	} else {
		_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)
		flashSession.AddFlash("Blog created", "message")
		flashSession.Save(r, w)
		http.Redirect(w, r, "/blog/list", http.StatusSeeOther)
	}

}

// DetailBlogGetHandler display blog detail
func DetailBlogGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	//fmt.Println(id)
	db := utils.GetDB()
	defer db.Close()
	var blog Blogs
	err := db.QueryRow("SELECT id, title, content, created_at, author as comments FROM blog WHERE id = ?", id).Scan(&blog.ID, &blog.Title, &blog.Content, &blog.CreatedAt, &blog.Author)
	page, _ := utils.GetPageStructure(w, r)
	idInt, _ := strconv.Atoi(id)
	blog.Comments = GetComments(idInt)
	page.PageData = blog
	err = tpl.ExecuteTemplate(w, "blog-detail.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetAuthor(userID int) Author {
	var author Author
	db := utils.GetDB()
	defer db.Close()
	err := db.QueryRow("SELECT id, fullname FROM `user` WHERE id = ?", userID).Scan(&author.ID, &author.Name)
	if err != nil {
		log.Fatalln(err)
	}

	return author
}

// EditBlogGetHandler edit form show for blog
func EditBlogGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LoginRequired(w, r)
	vars := mux.Vars(r)
	id := vars["id"]
	page, _ := utils.GetPageStructure(w, r)
	userID := utils.GetLogginUserID(r)
	db := utils.GetDB()
	defer db.Close()
	data := struct {
		ID      string
		Title   string
		Content template.HTML
	}{}
	err := db.QueryRow("SELECT id, title, content FROM blog where id = ? AND author = ?", id, userID).Scan(&data.ID, &data.Title, &data.Content)
	if err != nil {
		utils.RedirectWithMessage(w, r, "Blog not found or you are not authorized to perform this action!!!")
	} else {
		page.PageData = data
		err = tpl.ExecuteTemplate(w, "blog-edit.html", page)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

// EditBlogPostHandler update blog
func EditBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	utils.LoginRequired(w, r)
	r.ParseForm()
	vars := mux.Vars(r)
	id := vars["id"]
	userID := utils.GetLogginUserID(r)
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	db := utils.GetDB()
	defer db.Close()
	stmt, err := db.Prepare("UPDATE blog SET title = ?, content = ? WHERE id = ? AND author = ?")
	if err != nil {
		log.Fatalln(err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(title, content, id, userID)
	if err != nil {
		log.Fatalln(err)
	} else {
		_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)
		flashSession.AddFlash("Blog updated", "message")
		flashSession.Save(r, w)
		http.Redirect(w, r, "/blog/list", http.StatusSeeOther)
	}
}

// DeleteBlogHandler delete blog with passed id
func DeleteBlogHandler(w http.ResponseWriter, r *http.Request) {
	utils.LoginRequired(w, r)
	r.ParseForm()
	vars := mux.Vars(r)
	id := vars["id"]
	db := utils.GetDB()
	defer db.Close()
	_, err := db.Exec("DELETE FROM blog WHERE id = ? and author = ?", id, utils.GetLogginUserID(r))
	if err != nil {
		log.Fatalln(err)
	} else {
		msg := make(map[string]string)
		msg["msg"] = "Record Deleted"
		data, _ := json.Marshal(msg)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

// PopulateBlogTable blog table
func PopulateBlogTable(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	defer db.Close()
	for i := 5000; i < 20000; i++ {
		stmt, err := db.Prepare("INSERT INTO blog (title, content, created_at, status, author) VALUES(?,?,?,?,?)")
		if err != nil {
			log.Fatalln(err)
		}
		defer stmt.Close()
		content := "Hello Blog " + strconv.Itoa(i)
		_, err = stmt.Exec(content, content, time.Now().Format("2006-01-02 15:04:05"), 1, utils.GetLogginUserID(r))

	}
	_, _ = w.Write([]byte("Table populated"))
}

// AddBlogCommentPostHandler save posted comment
func AddBlogCommentPostHandler(w http.ResponseWriter, r *http.Request) {
	utils.LoginRequired(w, r)
	db := utils.GetDB()
	defer db.Close()
	r.ParseForm()
	blogID := r.PostFormValue("blog_id")
	content := r.PostFormValue("content")

	stmt, err := db.Prepare("INSERT INTO comment (blog_id, parent_id, content, created_at, created_by) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(blogID, 0, content, time.Now().Format("2006-01-02 15:04:05"), utils.GetLogginUserID(r))
	if err == nil {
		_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)
		flashSession.AddFlash("Comment added", "message")
		flashSession.Save(r, w)
		http.Redirect(w, r, "/blog/detail/"+blogID, http.StatusSeeOther)
	}
}
