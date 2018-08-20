package main

import (
	"Go-Webserver/webserver/articlehandler"
	"bufio"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	if _, noLog := os.Stat("log.txt"); os.IsNotExist(noLog) {
		newLog, err := os.Create("log.txt")
		if err != nil {
			log.Fatal(err)
		}
		newLog.Close()
	}
	dbString := readConfig("server.confi")
	var err error
	//var articleMux = http.NewServeMux()
	db, err := sql.Open("mysql", dbString)
	if err != nil {
		check(err)
	}
	//defer db.Close()
	err = db.Ping()
	if err != nil {
		check(err)
	}
	articlehandler.PassDataBase(db)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./src")))
	http.HandleFunc("/articles/", articlehandler.ReturnArticle)
	http.HandleFunc("/index.html", articlehandler.ReturnHomePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
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

func check(err error) {
	if err != nil {
		errorLog, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer errorLog.Close()
		log.SetOutput(errorLog)
		log.Println("An error has occured : ", err)
	}
}
