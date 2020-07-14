package user

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"

	utils "baseframe/utils/general"
	"baseframe/utils/pagination"
)

// User object
type User struct {
	ID       int
	Username string
	Password string
	Email    string
	Profile  Profile
}

// Profile for users
type Profile struct {
	ID        int
	Firstname string
	Lastname  string
	ContactNo string
	UserID    int
}

var perPage = 100
var tpl *template.Template

func init() {
	tpl = utils.GetTemplate()
}

// DashboardHandler show dashboard
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	utils.AdminLoginRequired(w, r)
	page, _ := utils.GetPageStructure(w, r)

	err := tpl.ExecuteTemplate(w, "dashboard.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// LoginGetHandler show form
func LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := utils.GetPageStructure(w, r)
	if page.IsLoggedIn {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
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
	// search only users with admin role
	err := db.QueryRow("SELECT id, username, password, email FROM `user` WHERE username = ? AND role_id = ?", username, 1).Scan(&dbUserID, &dbUsername, &dbPassword, &dbEmail)

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
			userSession.Values["isAdmin"] = true
			userSession.Save(r, w)
			http.Redirect(w, r, "/dashboard", 301)
		}

	} else {
		utils.RedirectWithMessage(w, r, "/login", "message", "Invalid username or password")
	}

}

// LogoutHandler logout user from session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_, userSession := utils.GetCookieStore(r, utils.STORE_SESSION)
	userSession.Values["isLoggedIn"] = false
	userSession.Values["isAdmin"] = false
	userSession.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ListGetHandler list all users
func ListGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.AdminLoginRequired(w, r)
	var profiles []Profile
	page, _ := utils.GetPageStructure(w, r)
	page.PageTitle = "All Users Listing"
	db := utils.GetDB()
	defer db.Close()

	q := "SELECT count(id) FROM `user` WHERE id > 1 ORDER BY id DESC"
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

	url := "/user/list"
	pager := pagination.New(totalRows, perPage, currentPage, url)
	page.Pager = pager

	c := make(chan Profile, 10)
	go getUsers(offset, c)
	for p := range c {
		profiles = append(profiles, p)
	}

	page.PageData = profiles
	err = tpl.ExecuteTemplate(w, "user-list.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

/*
func Ct(w http.ResponseWriter, r *http.Request) {
	var pro []Profile
	c := make(chan Profile, 10)
	go getUsers(1, c)
	for p := range c {
		pro = append(pro, p)
	}
	fmt.Println(pro)
	populateUsers()
}
*/
func getUsers(offset int, c chan Profile) {
	db := utils.GetDB()
	rows, err := db.Query("SELECT id, first_name, last_name, contact_no, user_id FROM `profile` WHERE id > 1 limit ?,?", offset, perPage)
	if err != nil {
		if err == sql.ErrNoRows {

		} else {
			log.Fatalln(err)
		}
	} else {
		defer db.Close()
		defer rows.Close()
		defer close(c)
		for rows.Next() {
			var id, userid int
			var firstname, lastname, contactno string
			err := rows.Scan(&id, &firstname, &lastname, &contactno, &userid)

			if err != nil {
				log.Fatalln(err)
			}
			c <- Profile{ID: id, Firstname: firstname, Lastname: lastname, ContactNo: contactno, UserID: userid}
		}
	}
}

func getUser(id int) (User, error) {
	var u User
	var p Profile
	db := utils.GetDB()
	defer db.Close()

	var userID, profileID int
	var firstname, lastname, contactno, username, email string

	err := db.QueryRow("SELECT p.id as profile_id, p.first_name, p.last_name, p.contact_no, u.id as user_id, u.username, u.email FROM `profile` p LEFT JOIN `user` u ON p.user_id = u.id WHERE u.id = ?", id).Scan(&profileID, &firstname, &lastname, &contactno, &userID, &username, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, errors.New("Profile not found")
		}
		return u, errors.New(err.Error())

	}

	p = Profile{Firstname: firstname, Lastname: lastname, ContactNo: contactno, ID: profileID}
	u = User{ID: userID, Username: username, Email: email, Profile: p}
	return u, nil

}

// CreateGetHandler create users
func CreateGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.AdminLoginRequired(w, r)
	page, _ := utils.GetPageStructure(w, r)
	page.PageTitle = "Create New User"
	err := tpl.ExecuteTemplate(w, "user-create.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// CreatePostHandler create users post request
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	utils.AdminLoginRequired(w, r)
	//page, _ := utils.GetPageStructure(w, r)

	db := utils.GetDB()
	defer db.Close()
	//_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)

	r.ParseForm()
	email := r.PostForm.Get("email")
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	firstname := r.PostForm.Get("firstname")
	lastname := r.PostForm.Get("lastname")
	contactno := r.PostForm.Get("contactno")

	if username == "" || password == "" || email == "" || firstname == "" || lastname == "" || contactno == "" {
		utils.RedirectWithMessage(w, r, "/user/create", "message", "All fields are required")
	}

	var existingUsername string
	err := db.QueryRow("SELECT username FROM `user` WHERE username='" + username + "' OR email = '" + email + "'").Scan(&existingUsername)

	if err != nil {
		if err == sql.ErrNoRows {
			hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			result, err := db.Exec("INSERT INTO `user` (username, password, email, role_id) VALUES (?,?,?,?)", username, hash, email, 2)
			userid, _ := result.LastInsertId()
			if err != nil {
				log.Fatal(err)
			} else {
				db.Exec("INSERT INTO profile (first_name, last_name, contact_no, user_id) VALUES(?,?,?,?)", firstname, lastname, contactno, userid)
				utils.RedirectWithMessage(w, r, "/user/create", "message", "Username "+existingUsername+" registered successfully")

			}
		} else {
			//log.Fatalln(err)
		}
	}
	if existingUsername != "" {
		utils.RedirectWithMessage(w, r, "/user/create", "message", "Username "+existingUsername+" already registered, please choose another")
	}

}

// EditGetHandler edit users
func EditGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.AdminLoginRequired(w, r)
	page, _ := utils.GetPageStructure(w, r)
	fmt.Printf("%+v\n", page)

	vars := mux.Vars(r)
	profileID, _ := strconv.Atoi(vars["id"])
	user, err := getUser(profileID)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Println(user)
	page.PageData = user
	page.PageTitle = "Update User \"" + user.Profile.Firstname + " " + user.Profile.Lastname + "\""
	err = tpl.ExecuteTemplate(w, "user-edit.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// EditPostHandler edit users post request
func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	utils.AdminLoginRequired(w, r)

	db := utils.GetDB()
	defer db.Close()
	//_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)

	r.ParseForm()
	email := r.PostForm.Get("email")
	username := r.PostForm.Get("username")
	firstname := r.PostForm.Get("firstname")
	lastname := r.PostForm.Get("lastname")
	contactno := r.PostForm.Get("contactno")
	id := r.PostForm.Get("id")

	if username == "" || email == "" || firstname == "" || lastname == "" || contactno == "" {
		utils.RedirectWithMessage(w, r, "/user/edit/"+id, "message", "All fields are required")
	}

	var existingUsername string
	err := db.QueryRow("SELECT username FROM `user` WHERE (username='"+username+"' OR email = '"+email+"') AND id != ?", id).Scan(&existingUsername)

	if err != nil {
		if err == sql.ErrNoRows {
			_, err := db.Exec("UPDATE `user` SET username = ?, email= ? WHERE id = ?", username, email, id)
			if err != nil {
				log.Fatal(err)
			} else {
				db.Exec("UPDATE profile SET first_name = ?, last_name =?, contact_no = ? WHERE user_id = ? ", firstname, lastname, contactno, id)
			}
			utils.RedirectWithMessage(w, r, "/user/edit/"+id, "message", "User updated successfully")
		}
	}
	if existingUsername != "" {
		utils.RedirectWithMessage(w, r, "/user/edit/"+id, "message", "Username "+existingUsername+" already registered, please choose another")
	}
}

// DeleteGetHandler delete users
func DeleteGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.AdminLoginRequired(w, r)
	page, _ := utils.GetPageStructure(w, r)

	err := tpl.ExecuteTemplate(w, "dashboard.html", page)
	if err != nil {
		log.Fatalln(err)
	}
}

func PopulateUsers(w http.ResponseWriter, r *http.Request) {

	//_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)
	vars := mux.Vars(r)
	start, _ := strconv.Atoi(vars["start"])
	end, _ := strconv.Atoi(vars["end"])

	db := utils.GetDB()
	defer db.Close()
	for i := start; i <= end; i++ { // always change the series

		//defer close(c)
		username := "test" + strconv.Itoa(i)
		email := username + "@localhost.com"
		password := "123456"

		var existingUsername string
		err := db.QueryRow("SELECT username FROM `user` WHERE username='" + username + "' OR email = '" + email + "'").Scan(&existingUsername)

		if err != nil {
			if err == sql.ErrNoRows {
				hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				res, err := db.Exec("INSERT INTO `user` (username, password, email, role_id) VALUES (?,?,?,?)", username, hash, email, 2)
				userid, _ := res.LastInsertId()
				firstname := "Test"
				lastname := strconv.Itoa(i)
				contactno := "99281" + strconv.Itoa(11111+rand.Intn(99999-11111))
				if err != nil {
					log.Println(err)
				} else {
					log.Println(email + " registered successfully")
					db.Exec("INSERT INTO profile (first_name, last_name, contact_no, user_id) VALUES(?,?,?,?)", firstname, lastname, contactno, userid)
				}
			} else {
				log.Println(err.Error())
			}
		}
		if existingUsername != "" {
			log.Println(email + " already registered")
		}

	}

}

func PopulateUserSecond(w http.ResponseWriter, r *http.Request) {

	//_, flashSession := utils.GetCookieStore(r, utils.FLASH_SESSION)
	vars := mux.Vars(r)
	start, _ := strconv.Atoi(vars["start"])
	end, _ := strconv.Atoi(vars["end"])

	db := utils.GetDB()
	defer db.Close()
	c := make(chan string)
	go func(c chan string) {
		defer close(c)

		for i := start; i <= end; i++ { // always change the series

			//defer close(c)
			username := "test" + strconv.Itoa(i)
			email := username + "@localhost.com"
			password := "123456"

			var existingUsername string
			err := db.QueryRow("SELECT username FROM `user` WHERE username='" + username + "' OR email = '" + email + "'").Scan(&existingUsername)

			if err != nil {
				if err == sql.ErrNoRows {
					hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
					res, err := db.Exec("INSERT INTO `user` (username, password, email, role_id) VALUES (?,?,?,?)", username, hash, email, 2)
					userid, _ := res.LastInsertId()
					firstname := "Test"
					lastname := strconv.Itoa(i)
					contactno := "99281" + strconv.Itoa(11111+rand.Intn(99999-11111))
					if err != nil {
						c <- err.Error()
					} else {
						c <- (email + " registered successfully")
						db.Exec("INSERT INTO profile (first_name, last_name, contact_no, user_id) VALUES(?,?,?,?)", firstname, lastname, contactno, userid)
					}
				} else {
					c <- err.Error()
				}
			}
			if existingUsername != "" {
				c <- (email + " already registered")
			}

		}
	}(c)
	for m := range c {
		fmt.Fprintln(w, m)
	}

}
