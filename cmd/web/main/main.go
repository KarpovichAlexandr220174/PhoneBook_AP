//package main
//
//import (
//	"log"
//	"net/http"
//	_ "time"
//)
//
//func main() {
//	mux := http.NewServeMux()
//
//
//	mux.HandleFunc("/", homeHandler)
//	mux.HandleFunc("/register", registerHandler)
//	mux.HandleFunc("/login", loginHandler)
//	mux.HandleFunc("/message", message)
//	mux.HandleFunc("/application", application)
//	mux.HandleFunc("/forgot", forgotPasswordHandler)
//
//	mux.HandleFunc("/city/", searchPageHandler)
//	mux.HandleFunc("/search/", searchHandler)
//
//	mux.HandleFunc("/city/hospitals/", searchHospitalsPageHandler)
//	mux.HandleFunc("/search/hospitals/", searchHospitalsHandler)
//
//	mux.HandleFunc("/city/schools/", searchSchoolsPageHandler)
//	mux.HandleFunc("/search/schools/", searchSchoolsHandler)
//
//	fileServer := http.FileServer(http.Dir("./ui/static/"))
//	mux.HandleFunc("/admin", adminHandler)
//	mux.HandleFunc("/admin/login", adminLoginHandler)
//	mux.HandleFunc("/admin/logout", adminLogoutHandler)
//	mux.HandleFunc("/admin/save-json", adminSaveJSONHandler)
//
//	//mux.HandleFunc("/admin/messages", adminSaveJSONHandler)
//	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
//	log.Println("Starting server on :5000")
//	err := http.ListenAndServe(":5000", mux)
//	log.Fatal(err)
//}
//

package main

import (
	_ "encoding/json"
	_ "html/template"
	_ "io/ioutil"
	"log"
	"net/http"
	"os"
	_ "strings"
	_ "time"

	"golang.org/x/time/rate"
)

var logger *log.Logger

func main() {
	initLogger()

	mux := http.NewServeMux()

	limiter := rate.NewLimiter(rate.Limit(1), 3)

	mux.HandleFunc("/", limitHandler(homeHandler, limiter))
	mux.HandleFunc("/register", limitHandler(registerHandler, limiter))
	mux.HandleFunc("/login", limitHandler(loginHandler, limiter))
	mux.HandleFunc("/message", limitHandler(message, limiter))
	mux.HandleFunc("/application", limitHandler(application, limiter))
	mux.HandleFunc("/forgot", limitHandler(forgotPasswordHandler, limiter))

	mux.HandleFunc("/city/", limitHandler(searchPageHandler, limiter))
	mux.HandleFunc("/search/", limitHandler(searchHandler, limiter))

	mux.HandleFunc("/city/hospitals/", limitHandler(searchHospitalsPageHandler, limiter))
	mux.HandleFunc("/search/hospitals/", limitHandler(searchHospitalsHandler, limiter))

	mux.HandleFunc("/city/schools/", limitHandler(searchSchoolsPageHandler, limiter))
	mux.HandleFunc("/search/schools/", limitHandler(searchSchoolsHandler, limiter))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.HandleFunc("/admin", limitHandler(adminHandler, limiter))
	mux.HandleFunc("/admin/login", limitHandler(adminLoginHandler, limiter))
	mux.HandleFunc("/admin/logout", limitHandler(adminLogoutHandler, limiter))
	mux.HandleFunc("/admin/save-json", limitHandler(adminSaveJSONHandler, limiter))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	logger.Println("Starting server on :5000")
	err := http.ListenAndServe(":5000", mux)
	if err != nil {
		logger.Fatal("Server failed to start:", err)
	}
}

func limitHandler(handler http.HandlerFunc, limiter *rate.Limiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			logger.Println("Rate limit exceeded")
			return
		}
		handler(w, r)
	}
}

func initLogger() {
	// Создаем папку logs, если она еще не существует
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatal("Не удалось создать каталог для логов:", err)
	}

	file, err := os.OpenFile("logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Не удалось открыть файл для логов:", err)
	}
	logger = log.New(file, "Server Log: ", log.Ldate|log.Ltime|log.Lshortfile)
}
