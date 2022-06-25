package main

import (
	"bookings/internals/config"
	"bookings/internals/driver"
	"bookings/internals/handlers"
	"bookings/internals/helpers"
	"bookings/internals/models"
	"bookings/internals/render"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

// Simple way without functions and methods
// func main() {

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		n, err := fmt.Fprintf(w, "Hello World!")
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		fmt.Println(n)
// 	})
// 	_ = http.ListenAndServe(":8000", nil)
// }
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// Organized way
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println("Starting web server...")
	// _ = http.ListenAndServe(":8000", nil)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {

	//What I am going to store in session?
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register([]models.Room{})
	gob.Register(models.Restriction{})

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction
	app.Session = session

	//Connect to Database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=Booking user=postgres password=")
	if err != nil {
		log.Fatal("Cannot connect to database...")
	}
	log.Println("Connected to database...")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Failed to create template cache", err)
		return nil, err
	}
	app.TemplateCache = tc
	app.AppCache = true
	render.NewRenderer(&app)

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)
	helpers.NewHelpers(&app)

	return db, nil
}
