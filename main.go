package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	var ERR error
	DB, ERR = sql.Open("sqlite3", "store.db")
	if ERR != nil {
		panic(ERR)
	}
	//defer DB.Close()
	URLTable := `CREATE TABLE IF NOT EXISTS url_cutter (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"url_id" TEXT,
		"url_redirect" TEXT);`
	Query, ERR := DB.Prepare(URLTable)
	if ERR != nil {
		panic(ERR)
	}
	Query.Exec()
}

// Функция для создания рандомной строки.
func rndGen(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// Функция получения url_id по ID(редеректа - url_redirect)
func DBGetFullURL(ID string) (error, string) {
	var s sql.NullString
	URLQuery := `SELECT "url_redirect" FROM "url_cutter" WHERE "url_id" = ?`
	error := DB.QueryRow(URLQuery, ID).Scan(&s)
	if error != nil {
		fmt.Println(error)
		return error, "none"
	}
	name := "Valued Customer"
	if s.Valid {
		name = s.String
	}
	return error, name
}

// Функция добавления нового редиректа в базу.
func DBAddURL(URL string) (error, string) {
	RNDString := rndGen(10)
	URLQuery := `INSERT INTO "url_cutter" ("url_redirect", url_id) VALUES(?,?)`
	_, error := DB.Exec(URLQuery, URL, RNDString)
	if error != nil {
		fmt.Println(error)
		return error, "none"
	}
	return error, RNDString
}

// Обработчик главной страницы.
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello world, from url_cutter!"))
}

// Обработчик для отображения содержимого редиректка.
func showCutter(w http.ResponseWriter, r *http.Request) {
	IDQuery := r.URL.Query().Get("id")
	ERR, URL := DBGetFullURL(IDQuery)
	if ERR != nil {
		http.NotFound(w, r)
		return
	}
	matched, _ := regexp.MatchString("^(http|https)://", URL)
	if !matched {
		URL = "http://" + URL
	}
	http.Redirect(w, r, URL, http.StatusSeeOther)
}

// Обработчик для создания нового редиректа.
func createCutter(w http.ResponseWriter, r *http.Request) {
	ERR, ID := DBAddURL(r.URL.Query().Get("url"))
	if ERR != nil {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("http://" + r.Host + "/url?id=" + ID))
	return
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/url", showCutter)
	mux.HandleFunc("/url/create", createCutter)
	error := http.ListenAndServe(":8080", mux)
	log.Fatal(error)
}
