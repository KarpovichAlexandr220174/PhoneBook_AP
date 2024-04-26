package main

import (
	"io/ioutil"
	"net/http"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка аутентификации администратора
	// В реальном приложении вам нужно будет использовать сессии или токены для проверки аутентификации
	isAuthenticated := checkAdminAuthentication(r)

	// Загрузка данных из JSON-файла (если аутентификация пройдена)
	var jsonData string
	if isAuthenticated {
		data, err := loadJSONData("phonebook.json")
		if err != nil {
			http.Error(w, "Error loading JSON data", http.StatusInternalServerError)
			return
		}
		jsonData = data
	}

	// Отображение HTML-шаблона
	tmpl.ExecuteTemplate(w, "admin.tmpl", struct {
		Authenticated bool
		JSONData      string
	}{
		Authenticated: isAuthenticated,
		JSONData:      jsonData,
	})
}

func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка учетных данных администратора и установка флага аутентификации (в данном примере флаг хранится в куки)
	// Замените это на вашу логику аутентификации
	username := r.FormValue("username")
	password := r.FormValue("password")
	isAuthenticated := authenticateAdmin(username, password)

	// Установка куки с флагом аутентификации
	if isAuthenticated {
		http.SetCookie(w, &http.Cookie{
			Name:  "admin_authenticated",
			Value: "true",
		})
	}

	// Редирект на страницу администратора
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func adminSaveJSONHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка аутентификации администратора перед сохранением JSON-данных
	if !checkAdminAuthentication(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получение данных из формы и сохранение в JSON-файле
	jsonData := r.FormValue("json_data")
	err := saveJSONData("phonebook.json", jsonData)
	if err != nil {
		http.Error(w, "Error saving JSON data", http.StatusInternalServerError)
		return
	}

	// Редирект на страницу администратора
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func adminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Очистка куки с флагом аутентификации
	http.SetCookie(w, &http.Cookie{
		Name:   "admin_authenticated",
		Value:  "",
		MaxAge: -1,
	})

	// Редирект на страницу администратора
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Функция для проверки аутентификации администратора
func checkAdminAuthentication(r *http.Request) bool {
	cookie, err := r.Cookie("admin_authenticated")
	if err != nil {
		return false
	}
	return cookie.Value == "true"
}

// Функция для аутентификации администратора (заглушка)
func authenticateAdmin(username, password string) bool {
	return username == "admin" && password == "55" // Замените на вашу логику аутентификации
}

// Функция для загрузки данных из JSON-файла
func loadJSONData(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Функция для сохранения данных в JSON-файле
func saveJSONData(filename, data string) error {
	return ioutil.WriteFile(filename, []byte(data), 0644)
}
