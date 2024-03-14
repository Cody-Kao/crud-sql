package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//var tmp *template.Template = template.Must(template.ParseFiles("templates/index.html"))

type Store struct {
	Storage map[string]*Book `json:"storage"`
}

type Book struct {
	Author string `json:"author"`
	Price  int    `json:"price"`
}

// in memory storage
var store Store = Store{map[string]*Book{"Lord Of Rings": {Author: "Tolkien", Price: 100},
	"The Witcher": {Author: "Andrzej", Price: 200}}}

func checkBook(bookName string) bool {
	if _, ok := store.Storage[bookName]; ok {
		return true
	}
	return false
}

func listAllBooks(w http.ResponseWriter) {
	// write this line to make browser treat them as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(store)
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
	store.Storage[bookName] = &Book{Author: author, Price: price}
	listAllBooks(w)
}

func Search(w http.ResponseWriter, r *http.Request) {
	bookName := mux.Vars(r)["bookName"]
	if !checkBook(bookName) {
		fmt.Fprintf(w, "book %s is not in store", bookName)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(store.Storage[bookName])
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
	store.Storage[bookName].Price = newPrice
	listAllBooks(w)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	bookName := mux.Vars(r)["bookName"]
	if !checkBook(bookName) {
		fmt.Fprintf(w, "book %s is not in store", bookName)
		return
	}
	delete(store.Storage, bookName)
	listAllBooks(w)
}
