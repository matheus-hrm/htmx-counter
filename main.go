package main

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Counter struct {
	value int
	mu    sync.Mutex
}

func (c *Counter) Add() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *Counter) Sub() {
	c.mu.Lock()
	c.value--
	c.mu.Unlock()
}

func (c *Counter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	counter := &Counter{}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("index.html")
		data := map[string]int {
			"CounterValue": counter.GetValue(),
		}
		tmpl.ExecuteTemplate(w, "index.html", data)																											
	})
	r.Post("/add", func(w http.ResponseWriter, _ *http.Request) {
		tmplStr := "<div id=\"counter\">{{.CounterValue}}</div>"
		tmpl := template.Must(template.New("counter").Parse(tmplStr))
		counter.Add()
		data := map[string]int {
			"CounterValue": counter.GetValue(),
		}
		tmpl.ExecuteTemplate(w, "counter", data)
	})
	r.Post("/sub", func(w http.ResponseWriter, _ *http.Request) {
		tmplStr := "<div id=\"counter\">{{.CounterValue}}</div>"
		tmpl := template.Must(template.New("counter").Parse(tmplStr))
		counter.Sub()
		data := map[string]int {
			"CounterValue": counter.GetValue(),
		}
		tmpl.ExecuteTemplate(w, "counter", data)
	})
	http.ListenAndServe("localhost:3000", r)
}
