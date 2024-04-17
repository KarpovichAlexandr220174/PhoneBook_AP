package main

import (
	"PhoneBook_AP/pkg/drivers"
	"log"
	"net/http"
	"time"
)

func main() {
	drivers.InitDB("user=postgres dbname=finalGo password=0000 sslmode=disable\n")
	mux := http.NewServeMux()

	http.Handle("/", TimeoutHandler(15*time.Second, mux))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/login", login)

	mux.HandleFunc("/application", application)

	mux.HandleFunc("/city/", searchPageHandler)
	mux.HandleFunc("/search/", searchHandler)

	mux.HandleFunc("/city/hospitals/", searchHospitalsPageHandler)
	mux.HandleFunc("/search/hospitals/", searchHospitalsHandler)

	mux.HandleFunc("/city/schools/", searchSchoolsPageHandler)
	mux.HandleFunc("/search/schools/", searchSchoolsHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.HandleFunc("/admin", adminHandler)
	mux.HandleFunc("/admin/login", adminLoginHandler)
	mux.HandleFunc("/admin/logout", adminLogoutHandler)
	mux.HandleFunc("/admin/save-json", adminSaveJSONHandler)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	log.Println("Starting server on :5000")
	err := http.ListenAndServe(":5000", mux)
	log.Fatal(err)
}
