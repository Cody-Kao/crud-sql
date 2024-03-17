package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/Cody-Kao/crud-sql/db"
	"github.com/gorilla/mux"
)

var DB *sql.DB

// query statement
// make sure to close sql.Stmt after using
var QueryOne *sql.Stmt
var QueryAll *sql.Stmt
var CreateRow *sql.Stmt
var UpdateRow *sql.Stmt
var DeleteRow *sql.Stmt

func SetDB(db *sql.DB) {
	DB = db
	// Initialize prepared statements
	QueryOne, _ = DB.Prepare("SELECT * FROM book WHERE name = $1")
	QueryAll, _ = DB.Prepare("SELECT * FROM book ORDER BY id")
	CreateRow, _ = DB.Prepare(`INSERT INTO book (name, author, price)
		 					VALUES($1, $2, $3) RETURNING id`)
	UpdateRow, _ = DB.Prepare("UPDATE book SET price = $1 WHERE name = $2 RETURNING id")
	DeleteRow, _ = DB.Prepare("DELETE FROM book WHERE name = $1 RETURNING id")
}

func checkBook(bookName string) bool {
	err := QueryOne.QueryRow(bookName).Scan()

	return err != nil
}

func listAllBooks(w http.ResponseWriter) {
	rows, err := QueryAll.Query()
	if err != nil {
		fmt.Fprint(w, "error occurs when fetch all books", err)
		return
	}

	var id int
	var name string
	var author string
	var price int
	var multiBook []db.Book
	for rows.Next() {
		err := rows.Scan(&id, &name, &author, &price)
		if err != nil {
			fmt.Fprint(w, "error occurs when scan book data", err)
			return
		}
		multiBook = append(multiBook, db.Book{ID: id, Name: name, Author: author, Price: price})
	}
	// write this line to make browser treat them as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(multiBook)
}

func ReadAll(w http.ResponseWriter, r *http.Request) {
	listAllBooks(w)
}

func Create(w http.ResponseWriter, r *http.Request) {
	bookName := r.URL.Query().Get("bookName")
	if checkBook(bookName) {
		fmt.Fprintf(w, "book %s is already in store", bookName)
		return
	}
	author := r.URL.Query().Get("author")
	price, err := strconv.Atoi(r.URL.Query().Get("price"))
	if err != nil {
		fmt.Fprintf(w, "error occurs: %s", err)
		return
	}
	var primaryKey int
	err = CreateRow.QueryRow(&bookName, &author, &price).Scan(&primaryKey)
	if err != nil {
		fmt.Fprint(w, "Error when create new data", err)
		return
	}
	fmt.Printf("book %s created with primary key of %d\n", bookName, primaryKey)
	listAllBooks(w)
}

func Search(w http.ResponseWriter, r *http.Request) {
	bookName := mux.Vars(r)["bookName"]
	if !checkBook(bookName) {
		fmt.Fprintf(w, "book %s is not in store", bookName)
		return
	}
	var id int
	var name string
	var author string
	var price int
	var oneBook db.Book
	err := QueryOne.QueryRow(bookName).Scan(&id, &name, &author, &price)
	if err != nil {
		fmt.Fprint(w, "error occurs when query one book", err)
		return
	}
	oneBook = db.Book{ID: id, Name: name, Author: author, Price: price}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(oneBook)
}

func Update(w http.ResponseWriter, r *http.Request) {
	bookName := r.URL.Query().Get("bookName")
	fmt.Println(bookName)
	if !checkBook(bookName) {
		fmt.Fprintf(w, "book %s is not in store", bookName)
		return
	}
	newPrice, err := strconv.Atoi(r.URL.Query().Get("Price"))
	if err != nil {
		fmt.Fprintf(w, "error occurs: %s", err)
		return
	}
	var primaryKey int
	err = UpdateRow.QueryRow(newPrice, bookName).Scan(&primaryKey)
	if err != nil {
		fmt.Fprint(w, "update error", err)
		return
	}
	fmt.Printf("%dth row is updated!\n", primaryKey)
	listAllBooks(w)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	bookName := mux.Vars(r)["bookName"]
	if !checkBook(bookName) {
		fmt.Fprintf(w, "book %s is not in store", bookName)
		return
	}
	var primaryKey int
	err := DeleteRow.QueryRow(bookName).Scan(&primaryKey)
	if err != nil {
		fmt.Fprint(w, "delete error", err)
		return
	}
	fmt.Printf("%dth row is deleted", primaryKey)
	listAllBooks(w)
}
