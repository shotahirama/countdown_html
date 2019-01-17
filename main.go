package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

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

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func htmlHandler(w http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFiles("./resources/index.html"))

	var s []map[string]string
	if Exists("countdown.yaml") {
		buf, err := ioutil.ReadFile("countdown.yaml")
		if err != nil {
			log.Fatalln(err)
		}
		var d CountDownData
		err = yaml.Unmarshal(buf, &d)
		if err != nil {
			log.Fatalln(err)
		}

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

func yamleditHandeler(w http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFiles("./resources/yamledit.html"))

	var s []map[string]string
	if Exists("countdown.yaml") {
		buf, err := ioutil.ReadFile("countdown.yaml")
		if err != nil {
			log.Fatalln(err)
		}
		var d CountDownData
		err = yaml.Unmarshal(buf, &d)
		for _, value := range d.CountDowns {
			sa := map[string]string{
				"name": value.Name,
				"date": value.Date,
			}
			s = append(s, sa)
		}
	}

	if err := t.ExecuteTemplate(w, "yamledit.html", s); err != nil {
		log.Fatal(err)
	}
}

type CountDowns []CountDown

func (a CountDowns) Len() int { return len(a) }
func (a CountDowns) Less(i, j int) bool {
	d1, _ := time.Parse("2006-01-02 15:04", a[i].Date)
	d2, _ := time.Parse("2006-01-02 15:04", a[j].Date)
	return d1.UnixNano() < d2.UnixNano()
}
func (a CountDowns) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func yamlpostHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	var cds CountDowns
	for key, value := range req.Form {
		// fmt.Println(key, value[0])
		cd := CountDown{key, value[0]}
		cds = append(cds, cd)
	}
	sort.Sort(CountDowns(cds))

	var d CountDownData
	d.CountDowns = cds

	buf, err := yaml.Marshal(&d)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("countdown.yaml", buf, 0644)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, req, "/yamledit", 301)
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	http.HandleFunc("/yamlpost", yamlpostHandler)
	http.HandleFunc("/yamledit", yamleditHandeler)
	http.HandleFunc("/countdown", countdownHandeler)
	http.HandleFunc("/", htmlHandler)

	http.ListenAndServe(":"+opts.Port, nil)
}
