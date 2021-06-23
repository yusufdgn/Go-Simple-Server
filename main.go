package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/tyler-sommer/stick"
)

type myHandler struct{}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: &myHandler{},
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = indexAction
	mux["/index"] = indexAction
	mux["/hello"] = helloAction
	fmt.Printf("listening - localhost:8080")
	server.ListenAndServe()
}

func indexAction(w http.ResponseWriter, r *http.Request) {
	j := map[string]stick.Value{"name": "Index : " + strconv.Itoa(rand.Intn(999999))}
	render("hello.html.twig", j, w, r)
}

func helloAction(w http.ResponseWriter, r *http.Request) {
	j := map[string]stick.Value{"name": "Hello : " + strconv.Itoa(rand.Intn(10000))}
	render("hello.html.twig", j, w, r)
}

func render(template string, j map[string]stick.Value, w http.ResponseWriter, r *http.Request) {
	var dir string
	var err error
	dir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir += "/template"

	env := stick.New(stick.NewFilesystemLoader(dir))
	env.Execute(template, w, j)
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	http.Error(w, "404 not found.", http.StatusNotFound)
	io.WriteString(w, "path: "+r.URL.String())
	return
}
