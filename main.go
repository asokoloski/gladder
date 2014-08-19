package main

import (
	"flag"
	"fmt"
	"github.com/bmizerany/pat"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/steveyen/gkvlite"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

func main() {
	db := flag.String("dbfile", os.Getenv("GLADDER_DB"), "path to use for gladder database (defaults to value of $GLADDER_DB)")
	httpAddr := flag.String("http", os.Getenv("GLADDER_HTTP_ADDR"), "http port (':80' for example. defaults to value of $GLADDER_HTTP_ADDR)")
	if *db == "" {
		*db = "gladder-db.gkv"
	}
	flag.Parse()
	dbFile, err := os.OpenFile(*db, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0x700)
	if err != nil {
		log.Fatal(err)
	}
	store, err := gkvlite.NewStore(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	gladder := NewGladder(store)
	gw := &GladderWeb{gladder}

	log.Println("Starting up")
	log.Println("Listening on", *httpAddr)
	log.Println("Using database file", *db)
	log.Println("Assets:")
	for _, name := range AssetNames() {
		log.Println("  ", name)
	}
	log.Fatal(http.ListenAndServe(*httpAddr, gw.Mux()))
}

type GladderWeb struct {
	*Gladder
}

func (w *GladderWeb) Mux() http.Handler {
	mux := pat.New()
	mux.Get("/resources/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, ""}))
	mux.Get("/create_player", http.HandlerFunc(w.createPlayer))
	mux.Post("/create_player", http.HandlerFunc(w.createPlayer))
	mux.Post("/player/:name/", http.HandlerFunc(w.editPlayer))
	mux.Get("/", http.HandlerFunc(w.index))
	return mux
}

func (gw *GladderWeb) index(w http.ResponseWriter, r *http.Request) {
	users, err := gw.Gladder.GetUsers()
	sort.Sort(users)
	if err != nil {
		panic(err)
	}
	render(w, "resources/templates/index.html", map[string]interface{}{
		"Users": users,
	})
}

func (gw *GladderWeb) editPlayer(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	uname := r.URL.Query().Get(":name")
	user, err := gw.Gladder.GetUser(uname)
	if err != nil {
		log.Println(err)
		return
	}
	rankStr := r.Form.Get("rank")
	rank, err := strconv.Atoi(rankStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	user.Rank = Rank(rank)
	fmt.Println("Setting rank of", uname, "to", rank)
	gw.Gladder.SaveUser(user)
	http.Redirect(w, r, "/", 302)
}

func (gw *GladderWeb) createPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			return
		}
		uname := r.Form.Get("username")
		if uname != "" {
			existing, err := gw.Gladder.GetUser(uname)
			if err != nil {
				log.Println("user already exists:", uname)
				return
			}
			rankStr := r.Form.Get("rank")
			rank, err := strconv.Atoi(rankStr)
			if err != nil {
				fmt.Println(err)
				return
			}
			if existing != nil {
				err = gw.Gladder.CreateUser(uname, Rank(rank))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
		http.Redirect(w, r, "/", 302)
		return
	}
	render(w, "resources/templates/create_user.html", map[string]interface{}{})
}

func render(w http.ResponseWriter, templateFile string, context interface{}) {
	var err error
	var tmpl *template.Template
	file, err := Asset(templateFile)
	if err != nil {
		goto ERR
	}
	tmpl, err = template.New(templateFile).Parse(string(file))
	if err != nil {
		goto ERR
	}
	err = tmpl.Execute(w, context)
	if err != nil {
		goto ERR
	}
	return
ERR:
	log.Println(err)
	http.Error(w, fmt.Sprintf("%s", err), 500)
	return
}
