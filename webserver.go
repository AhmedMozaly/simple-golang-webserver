package main

import (
	"bufio"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//var templates *template.Template
type postBin struct {
	Posters []article
}

type article struct {
	ID       int
	Title    string
	PostText string
	Date     string
	ImageURL string
	Tags     []string
}

var db *sql.DB

func main() {
	dbString := readConfig("server.confi")
	var err error
	//var articleMux = http.NewServeMux()
	db, err = sql.Open("mysql", dbString)
	if err != nil {
		log.Fatal("dbConnection failed : ", err)
	}
	//defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal("dbConnection failed : ", err)
	}
	http.Handle("/", http.FileServer(http.Dir("./src")))
	http.HandleFunc("/articles/", articleHandler)
	http.HandleFunc("/index.html", homePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	requestURI := strings.SplitAfter(r.RequestURI, "/")
	articleID, err := strconv.Atoi(requestURI[len(requestURI)-1])
	template, err := template.ParseFiles("./src/articles/article.html")
	article := getArticle(articleID)

	err = template.Execute(w, article)
	if err != nil {
		//go fileNotFound(err, w, r)
		log.Fatal("Fatal parsing error : ", err)
	}

}

func homePage(w http.ResponseWriter, r *http.Request) {

	var allp = frontPagePosts(db)
	//var posts = postBin{testPost}
	templates, err := template.ParseFiles("./src/index.html")
	if err != nil {
		log.Fatal("Parsing error : ", err)
	}
	homePage := templates.Lookup("index.html")
	homePage.Execute(w, allp)

}

func frontPagePosts(db *sql.DB) postBin {
	returnPost, err := db.Query("SELECT * FROM blog.blog_posts LIMIT 5")
	if err != nil {
		log.Fatal("DB statement failed : ", err)
	}
	var dbResults []article
	defer returnPost.Close()
	for returnPost.Next() {
		var bp article
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
	returnArticle := s.QueryRow(id)
	returnArticle.Scan(&ar.ID, &ar.Title, &ar.PostText, &ar.Date)
	return ar
}

func readConfig(s string) string {
	config, err := os.Open(s)
	if err != nil {
		log.Fatal("File open failure : ", err)
	}
	defer config.Close()

	scanner := bufio.NewScanner(config)
	scanner.Scan()
	return scanner.Text()

}

/*func check(err error) {
	if err != nil {
		if _, noLog := os.Stat("log.txt"); os.IsNotExist(noLog) {
			newLog, err := os.Create("log.txt")
			if err != nil {
				log.Fatal(err)
			}
			newLog.Close()
		}
		errorLog, err := os.Open("log.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer errorLog.Close()
		log.SetOutput(errorLog)
		log.Printf("An error has occured : ", err)
	}
}*/
