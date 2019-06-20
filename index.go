package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

type user struct {
	Id      int
	Name    string
	Country string
}

type users struct {
	Titre     string
	UsersList []user
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func accueil(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./test.db")
	checkErr(err)

	rows, err := db.Query("SELECT * FROM users")
	checkErr(err)
	var id int
	var name string
	var country string
	var liste []user

	for rows.Next() {
		err = rows.Scan(&id, &name, &country)
		checkErr(err)
		liste = append(liste, user{Id: id, Name: name, Country: country})
	}
	rows.Close()
	db.Close()
	t, _ := template.ParseFiles("index.html")
	data := users{Titre: "Mes copains", UsersList: liste}
	fmt.Println(data)
	t.Execute(w, data)
}

func ajouter(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./test.db")
	checkErr(err)

	r.ParseForm()

	nom := r.Form["nom"][0]
	pays := r.Form["pays"][0]

	stmt, err := db.Prepare("INSERT INTO users(name, country) VALUES(?,?)")
	checkErr(err)

	stmt.Exec(nom, pays)

	db.Close()
	http.Redirect(w, r, "/", http.StatusFound)
}

func modifier(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("modif.html")

	r.ParseForm()

	id, _ := strconv.Atoi(r.Form["idForm"][0])
	nom := r.Form["nom"][0]
	pays := r.Form["pays"][0]

	data := user{Id: id, Name: nom, Country: pays}

	t.Execute(w, data)
}

func executerModif(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./test.db")
	checkErr(err)

	r.ParseForm()

	id, _ := strconv.Atoi(r.Form["idForm"][0])
	nom := r.Form["nom"][0]
	pays := r.Form["pays"][0]

	stmt, err := db.Prepare("UPDATE users SET name=?, country=? WHERE id=?")
	checkErr(err)

	stmt.Exec(nom, pays, id)

	db.Close()

	http.Redirect(w, r, "/", http.StatusFound)
}

func supprimer(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./test.db")
	checkErr(err)

	r.ParseForm()
	IdForm := r.Form["idForm"][0]

	stmt, err := db.Prepare("DELETE FROM users WHERE id=?")
	checkErr(err)

	stmt.Exec(IdForm)

	db.Close()
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", accueil)
	http.HandleFunc("/ajouter/", ajouter)
	http.HandleFunc("/supprimer/", supprimer)
	http.HandleFunc("/modifier/", modifier)
	http.HandleFunc("/executerModif/", executerModif)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
