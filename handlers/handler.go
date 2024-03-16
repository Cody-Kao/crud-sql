package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Cody-Kao/crud-sql/db"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var DB *gorm.DB
var oneBook db.Book
var multiBook []db.Book

func SetDB(db *gorm.DB) {
	DB = db
}

func checkBook(bookName string) bool {
	res := DB.Where("name = ?", bookName).First(&oneBook)

	return res.Error == nil
}

func listAllBooks(w http.ResponseWriter) {
	res := DB.Order("id asc").Find(&multiBook)
	// write this line to make browser treat them as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(multiBook)
	fmt.Println("selected rows:", res.RowsAffected)
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

	res := DB.Create(&db.Book{Name: bookName, Author: author, Price: price})
	if res.Error != nil {
		fmt.Fprint(w, "Error when create new data", err)
		return
	}
	listAllBooks(w)
}

func Search(w http.ResponseWriter, r *http.Request) {
	bookName := mux.Vars(r)["bookName"]
	if !checkBook(bookName) {
		fmt.Fprintf(w, "book %s is not in store", bookName)
		return
	}
	DB.Where("name = ?", bookName).First(&oneBook)
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
	res := DB.Model(&oneBook).Where("name = ?", bookName).Update("price", newPrice)
	if res.Error != nil {
		fmt.Fprint(w, "update error")
		return
	}

	listAllBooks(w)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	bookName := mux.Vars(r)["bookName"]
	if !checkBook(bookName) {
		fmt.Fprintf(w, "book %s is not in store", bookName)
		return
	}
	res := DB.Where("name = ?", bookName).Delete(&oneBook)
	if res.Error != nil {
		fmt.Fprint(w, "delete error")
		return
	}
	listAllBooks(w)
}
