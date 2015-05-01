package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"
	"strconv"

	"./data"
)

var (
	d         data.Provider
	port      = flag.Int("port", 8080, "HTTP port.")
	tmplDir   = flag.String("template_dir", "tmpl", "Directory to load templates from.")
	staticDir = flag.String("static_dir", "static", "Directory to load static files from. Served with /static/ URI prefix.")
	dataFile  = flag.String("data", "", "File to use for data storage.")
	menuTmpl  *template.Template
)

func loadTemplates(dir string) error {
	var err error
	intToDay := map[Weekday]string{
		Mon: "Понеділок",
		Tue: "Вівторок",
		Wed: "Середа",
		Thu: "Четвер",
		Fri: "П’ятниця",
		Sat: "Субкота",
		Sun: "Неділя",
	}
	funcMap := template.FuncMap{
		"dayname": func(d Weekday) (string, error) {
			if r, found := intToDay[d]; found {
				return r, nil
			} else {
				return "", errors.New("unknown day of the week")
			}
		},
	}
	menuTmpl, err = template.New("menu.tmpl").Funcs(funcMap).ParseFiles(
		filepath.Join(dir, "menu.tmpl"),
		filepath.Join(dir, "head.tmpl"),
		filepath.Join(dir, "footer.tmpl"),
		filepath.Join(dir, "breakfast.tmpl"),
		filepath.Join(dir, "lunch.tmpl"),
		filepath.Join(dir, "dinner.tmpl"),
		filepath.Join(dir, "dish.tmpl"),
		filepath.Join(dir, "data.tmpl"))
	return err
}

func makeMenu(dishes []data.Dish) (WeekMenu, error) {
	byType := makeDishMap(dishes)

	menu := WeekMenu{}
	for day := Mon; day <= Sun; day++ {
		menu[day] = DayMenu{
			Breakfast: genBreakfast(byType),
			Lunch:     genLunch(byType, day),
			Dinner:    genDinner(byType),
		}
	}

	return menu, nil
}

func generateMenu(w http.ResponseWriter, req *http.Request) {
	if d == nil {
		log.Print("d == nil")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	dishes, err := d.Dishes()
	if err != nil {
		log.Printf("Failed to get dishes: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	menu, err := makeMenu(dishes)
	if err != nil {
		log.Printf("Failed to generate a menu: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := menuTmpl.Execute(w, menu); err != nil {
		log.Printf("Failed to render menu: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func generateBreakfast(w http.ResponseWriter, req *http.Request) {
	if d == nil {
		log.Print("d == nil")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	dishes, err := d.Dishes()
	if err != nil {
		log.Printf("Failed to get dishes: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := menuTmpl.ExecuteTemplate(w, "breakfast.tmpl", genBreakfast(makeDishMap(dishes))); err != nil {
		log.Printf("Failed to render menu: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func generateLunch(w http.ResponseWriter, req *http.Request) {
	if d == nil {
		log.Print("d == nil")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	dishes, err := d.Dishes()
	if err != nil {
		log.Printf("Failed to get dishes: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	day, err := strconv.Atoi(req.FormValue("day"))
	if err != nil {
		log.Printf("Failed to get day: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := menuTmpl.ExecuteTemplate(w, "lunch.tmpl", genLunch(makeDishMap(dishes),Weekday(day))); err != nil {
		log.Print("Failed to render menu: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func generateDinner(w http.ResponseWriter, req *http.Request) {
	if d == nil {
		log.Print("d == nil")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	dishes, err := d.Dishes()
	if err != nil {
		log.Printf("Failed to get dishes: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := menuTmpl.ExecuteTemplate(w, "dinner.tmpl", genDinner(makeDishMap(dishes))); err != nil {
		log.Printf("Failed to render menu: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func dump(w http.ResponseWriter, req *http.Request) {
	if d == nil {
		log.Print("d == nil")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	dump, err := data.Dump(d)
	if err != nil {
		log.Printf("Failed to dump data: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(dump); err != nil {
		log.Printf("Failed to write response: %s", err)
	}
}

func restore(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("/restore invoked with %q", req.Method)
		http.Error(w, "Should be only invoked with POST", http.StatusBadRequest)
		return
	}
	if d == nil {
		log.Print("d == nil")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		log.Printf("Failed to parse mulipart form: %s", err)
		http.Error(w, "Failed to parse the request", http.StatusBadRequest)
		return
	}
	if req.MultipartForm == nil {
		log.Print("req.MultipartForm == nil")
		http.Error(w, "No multipart form", http.StatusBadRequest)
		return
	}
	if req.MultipartForm.File == nil || req.MultipartForm.File["dump"] == nil || len(req.MultipartForm.File["dump"]) != 1 {
		log.Printf("No dump file detected in\n%#v", req.MultipartForm)
		http.Error(w, "Malformed form", http.StatusBadRequest)
		return
	}
	f, err := req.MultipartForm.File["dump"][0].Open()
	if err != nil {
		log.Printf("Failed to open dump file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	b, err := ioutil.ReadAll(f)
	f.Close()
	if err != nil {
		log.Printf("Failed to read dump file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := d.Restore(b); err != nil {
		log.Printf("Failed to restore from dump: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func dataUI(w http.ResponseWriter, req *http.Request) {
	if err := menuTmpl.ExecuteTemplate(w, "data.tmpl", nil); err != nil {
		log.Printf("Failed to render data UI: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	flag.Parse()
	if *dataFile == "" {
		log.Fatal("Must provide --data.")
	}
	var err error
	d, err = data.NewLocalProvider(*dataFile)
	if err != nil {
		log.Fatal("NewLocalProvider: ", err)
	}

	if err := loadTemplates(*tmplDir); err != nil {
		log.Fatal("loadTemplates: ", err)
	}

	seed := time.Now().UnixNano()
	log.Printf("Initializing random number generator with seed %d", seed)
	rand.Seed(seed)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(*staticDir))))
	http.HandleFunc("/", generateMenu)
	http.HandleFunc("/gen/breakfast", generateBreakfast)
	http.HandleFunc("/gen/lunch", generateLunch)
	http.HandleFunc("/gen/dinner", generateDinner)
	http.HandleFunc("/dump", dump)
	http.HandleFunc("/restore", restore)
	http.HandleFunc("/data", dataUI)
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting server at %q...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListerAndServe: ", err)
	}
}
