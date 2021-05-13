package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Person struct {
	Name string
}

func getPeople() []Person {
	return []Person{
		{Name: "Joe"},
		{Name: "Bob"},
	}
}

func getPerson(id int) (Person, error) {
	people := getPeople()

	if id < 0 || id >= len(people) {
		return Person{}, fmt.Errorf("Could not find person at index: %d", id)
	}

	return people[id], nil

}

func TurboFrame(id string) func(w http.ResponseWriter, f func(http.ResponseWriter)) {
	return func(w http.ResponseWriter, f func(http.ResponseWriter)) {
		open := fmt.Sprintf("<turbo-frame id=\"%s\">", id)
		io.WriteString(w, open)
		f(w)
		io.WriteString(w, "</turbo-frame>")
	}
}

func HandleGetPeople(w http.ResponseWriter, r *http.Request) {
	people := getPeople()

	rend := func(w http.ResponseWriter) {
		if err := renderPeople(r.Context(), w, people); err != nil {
			log.Println("error", err)
		}
	}

	//rend(w)
	TurboFrame("people")(w, rend)
}

func HandleGetPerson(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Println("error", err)
		http.NotFound(w, r)
		return
	}

	p, err := getPerson(id)

	if err != nil {
		log.Println("error", err)
		http.NotFound(w, r)
		return
	}

	if err := renderPerson(r.Context(), w, p); err != nil {
		log.Println("error", err)
	}
}
func htmlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

func main() {
	var port string

	flag.StringVar(&port, "port", "80", "port to run the website on")
	flag.Parse()

	r := mux.NewRouter()
	r.Use(htmlMiddleware)

	r.HandleFunc("/people", HandleGetPeople)
	r.HandleFunc("/person/{id}", HandleGetPerson)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	http.Handle("/", r)
	fmt.Println("Listening on port:", port)
	http.ListenAndServe(":"+port, nil)
}
