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

// Функция для создания рандомной строки.
func rnd_gen(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	fmt.Println(b.String())
	return b.String()
}

// Функция для создния таблицы.
func db_create() {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	url_table := `CREATE TABLE IF NOT EXISTS url_cutter (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "url_id" TEXT,
        "url_redirect" TEXT);`
	query, err := db.Prepare(url_table)
	if err != nil {
		panic(err)
	}
	query.Exec()
}

// Функция получения url_id по ID(редеректа - url_redirect)
func db_get_full_url(url_id string) (error, string) {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var s sql.NullString
	url_table := `SELECT "url_redirect" FROM "url_cutter" WHERE "url_id" = ?`
	error := db.QueryRow(url_table, url_id).Scan(&s)
	if error != nil {
		fmt.Println(error)
		return error, "none"
	}
	name := "Valued Customer"
	if s.Valid {
		name = s.String
	}
	fmt.Println(name)
	return error, name
}

// Функция добавления нового редиректа в базу.
func db_add_url(url string) (error, string) {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rnd_id := rnd_gen(10)
	url_table := `INSERT INTO "url_cutter" ("url_redirect", url_id) VALUES(?,?)`
	_, error := db.Exec(url_table, url, rnd_id)
	if error != nil {
		fmt.Println(error)
		return error, "none"
	}
	return error, rnd_id
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
	url_id_string := r.URL.Query().Get("id")
	err, url := db_get_full_url(url_id_string)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	matched, _ := regexp.MatchString("^(http|https)://", url)
	if !matched {
		url = "http://" + url
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// Обработчик для создания нового редиректа.
func createCutter(w http.ResponseWriter, r *http.Request) {
	err, id := db_add_url(r.URL.Query().Get("url"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("http://" + r.Host + "/url?id=" + id))
	return
}

func main() {
	go db_create()
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/url", showCutter)
	mux.HandleFunc("/url/create", createCutter)

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
