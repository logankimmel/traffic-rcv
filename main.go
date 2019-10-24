package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type (
	routeSet []route
	route    struct {
		Name        string
		Method      string
		Patter      string
		HandlerFunc http.HandlerFunc
	}
)

var routes = routeSet{
	route{"Home", "GET", "/", home},
	route{"Bad", "GET", "/test", bad},
}

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s\t",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Patter).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "OK")
}

func bad(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	fmt.Fprint(w, "ERROR")
}

func main() {
	router := newRouter()
	log.Println("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
