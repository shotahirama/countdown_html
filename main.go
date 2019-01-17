package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	flags "github.com/jessevdk/go-flags"
	yaml "gopkg.in/yaml.v2"
)

type CountDownData struct {
	CountDowns []CountDown `yaml:"countdowns"`
}

type CountDown struct {
	Name string `yaml:"name"`
	Date string `yaml:"date"`
}

type Options struct {
	Port string `short:"p" long:"port" description:"port" default:"12000"`
}

var opts Options

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
	for _, value := range d.CountDowns {
		v := url.Values{}
		v.Set("name", value.Name)
		v.Set("timelimit", value.Date)
		sa := map[string]string{
			"name": value.Name,
			"URL":  "/countdown" + "?" + v.Encode(),
		}
		s = append(s, sa)
	}

	if err := t.ExecuteTemplate(w, "index.html", s); err != nil {
		log.Fatal(err)
	}
}

func countdownHandeler(w http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFiles("./resources/base.html"))

	if err := t.ExecuteTemplate(w, "base.html", "countdown"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	http.HandleFunc("/countdown", countdownHandeler)
	http.HandleFunc("/", htmlHandler)

	http.ListenAndServe(":"+opts.Port, nil)
}
