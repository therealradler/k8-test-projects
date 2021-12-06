package main

import (
	"html/template"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var logFields = log.Fields{
	"application": "go-wiki",
}
var port = "8081"

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		errorLogger(err.Error())
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	requestLogger(r)
	p, err := loadPage(title)
	if err != nil {
		errorLogger(err.Error())
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, "view", p)
}
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	requestLogger(r)
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	requestLogger(r)
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		errorLogger(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	log.WithFields(logFields).WithField("event", "render template").Info("returned template: ", tmpl)
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		errorLogger(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func requestLogger(r *http.Request) {
	log.WithFields(logFields).WithFields(log.Fields{
		"uri":    r.RequestURI,
		"method": r.Method,
	}).Info("request completed")
}

func errorLogger(err string) {
	log.WithFields(logFields).WithField("event", "error").Error(err)
}

func main() {
	log.WithFields(logFields).Info("Server Start")
	serverPort := ":" + port
	log.WithFields(logFields).Info(serverPort)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.WithFields(logFields).WithField("event", "start server").Info("Starting Server on Port: ", port)
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.WithFields(logFields).WithField("event", "start server").Fatal(err)
	}

}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}
