package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	yaml "gopkg.in/yaml.v2"
)

type CountDownData struct {
	CountDowns []CountDown `yaml:"countdowns"`
}

type CountDown struct {
	Name string `yaml:"name"`
	Date string `yaml:"date"`
}

var r *mux.Router

func htmlHandler(w http.ResponseWriter, req *http.Request) {

	buf, err := ioutil.ReadFile("countdown.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	var d CountDownData
	err = yaml.Unmarshal(buf, &d)
	if err != nil {
		log.Fatalln(err)
	}
	t := template.Must(template.ParseFiles("./resources/index.html"))

	var s []map[string]string
	rtr := mux.NewRouter()
	for _, value := range d.CountDowns {
		rtr.Path("/"+value.Name).Queries("name", "{name}", "timelimit", "{timelimit}").HandlerFunc(hogehandle).Name(value.Name)
		url, err := rtr.Get(value.Name).URL("name", value.Name, "timelimit", value.Date)
		if err == nil {
			r.Handle("/"+value.Name, rtr)
			sa := map[string]string{
				"name": value.Name,
				"URL":  url.String(),
			}
			s = append(s, sa)
		}
	}

	if err := t.ExecuteTemplate(w, "index.html", s); err != nil {
		log.Fatal(err)
	}
}

func hogehandle(w http.ResponseWriter, r *http.Request) {

	key := r.FormValue("name")
	timelimit := r.FormValue("timelimit")

	t := template.Must(template.ParseFiles("./resources/base.html"))

	m := map[string]string{
		"key1": key,
		"key2": timelimit,
	}

	if err := t.ExecuteTemplate(w, "base.html", m); err != nil {
		log.Fatal(err)
	}
}

func main() {
	buf, err := ioutil.ReadFile("countdown.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	var d CountDownData
	err = yaml.Unmarshal(buf, &d)
	if err != nil {
		log.Fatalln(err)
	}

	r = mux.NewRouter()
	r.PathPrefix("/resources/").Handler(http.StripPrefix("/resources", http.FileServer(http.Dir("."+"/resources/"))))
	r.HandleFunc("/", htmlHandler)
	http.Handle("/", r)

	http.ListenAndServe(":12000", nil)
}
