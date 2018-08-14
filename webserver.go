package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//var templates *template.Template
type postBin struct {
	Posters []splashPost
}

type article struct {
	ID       int
	Title    string
	PostText string
	Date     string
	ImageURL string
	Tags     []string
}

type splashPost struct {
	ID       int
	Title    string
	PostText string
	Date     string
}

var db *sql.DB

func main() {
	var err error
	var articleMux = http.NewServeMux()
	db, err = sql.Open("mysql", "")
	if err != nil {
		log.Fatal("dbConnection failed : ", err)
	}
	//defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal("dbConnection failed : ", err)
	}

	http.Handle("/", http.FileServer(http.Dir("./src")))

	http.HandleFunc("/index.html", homePage)
	//articleMux.Handle("/articles/", nil)
	articleMux.HandleFunc("/", articleHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	requestURI := strings.SplitAfter(r.RequestURI, "/")
	articleID, err := strconv.Atoi(requestURI[len(requestURI)-1])
	template, err := template.ParseFiles("./src/articles/articles.html")
	article := getArticle(articleID)
	template.Execute(w, article)
	if err != nil {
		go fileNotFound(err, w, r)
	}

}

func homePage(w http.ResponseWriter, r *http.Request) {

	var allp = frontPagePosts(db)
	//var posts = postBin{testPost}
	templates, err := template.ParseFiles("./src/index.html")
	if err != nil {
		fileNotFound(err, w, r)
	}
	homePage := templates.Lookup("index.html")
	homePage.Execute(w, allp)

}

func frontPagePosts(db *sql.DB) postBin {
	returnPost, err := db.Query("SELECT * FROM blog.blog_posts LIMIT 3")
	if err != nil {
		log.Fatal("DB statement failed : ", err)
	}
	var dbResults []splashPost
	defer returnPost.Close()
	for returnPost.Next() {
		var bp splashPost
		returnPost.Scan(&bp.ID, &bp.Title, &bp.PostText, &bp.Date)
		bp.PostText = bp.PostText[0:40]
		dbResults = append(dbResults, bp)
	}
	pb := postBin{
		dbResults,
	}
	return pb
}

func getArticle(id int) article {
	var ar = article{}
	s, err := db.Prepare("SELECT * from blog.blog_posts WHERE idblog_posts = ?")
	if err != nil {
		log.Fatal("Statement prep failed : ", err)
	}
	returnArticle, err := s.Query(id)
	defer returnArticle.Close()
	returnArticle.Scan(&ar.ID, &ar.Title, &ar.PostText, &ar.Date)
	return ar
}

func fileNotFound(err error, w http.ResponseWriter, r *http.Request) {
	log.Print("File not found, ", err)
	http.ServeFile(w, r, "./src/404.html")
}
