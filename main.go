package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"fmt"
	"github.com/speps/go-hashids"
	"time"

	"github.com/gorilla/mux"
	"net/http"
)

var db *sql.DB
var err error

func main()  {

	db, err = sql.Open("sqlite3", "urls.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/create/{long}", create)
	router.HandleFunc("/go/{short}", access)
	http.ListenAndServe(":9000", router)
}


func access(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	rows, err := db.Query("select * from foo where shortUrl='" + params["short"] + "'")
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var shortUrl, longUrl string
		rows.Scan(&shortUrl, &longUrl)
		http.Redirect(w, r, "http://" +  longUrl, 302)
	}
}

func create(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	rows, err := db.Query("select * from foo where longUrl='" + params["long"] + "'")
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var shortUrl, longUrl string
		rows.Scan(&shortUrl, &longUrl)
		w.Write([]byte(shortUrl))
		return
	}

	short := addIntoDB(params["long"])

	w.Write([]byte(short))

}




func hash(toHash string) (hashedValue string){
	hd := hashids.NewData()
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	hashedValue, _ = h.Encode([]int{int(now.Unix())})
	return
}

func addIntoDB(long string) (short string) {

	short = hash(long)

	rows, _ := db.Query("select * from db where shortUrl='" + short +"'")
	if rows != nil {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}

	stmt, err := tx.Prepare("insert into foo(shortUrl, longUrl) values(?, ?)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(short, long)
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()
	return
}