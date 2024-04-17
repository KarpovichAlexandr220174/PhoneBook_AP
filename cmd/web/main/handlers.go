package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type PhonebookData struct {
	City      string
	Phonebook []PhoneEntry
	Hospitals []HospitalEntry
	Schools   []SchoolEntry
}

type PhoneEntry struct {
	Phone string `json:"phone"`
	Name  string `json:"name"`
}

type HospitalEntry struct {
	Name    string `json:"name"`
	Number  string `json:"number"`
	Address string `json:"address"`
}

type SchoolEntry struct {
	Name    string `json:"name"`
	Number  string `json:"number"`
	Address string `json:"address"`
}

//---------------------------------------------------------------------------------

var tmpl = template.Must(template.ParseGlob("ui/html/*.tmpl"))

func application(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "main-page.tmpl", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.tmpl", nil)
}

func searchSchoolsPageHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Path[len("/city/schools/"):]
	if city == "" {
		http.NotFound(w, r)
		return
	}
	tmpl.ExecuteTemplate(w, "search-schools.tmpl", PhonebookData{City: city})
}

func searchSchoolsHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Path[len("/search/schools/"):]
	if city == "" {
		http.NotFound(w, r)
		return
	}

	// Чтение данных из JSON
	file, err := ioutil.ReadFile("phonebook.json")
	if err != nil {
		http.Error(w, "Error reading phonebook data", http.StatusInternalServerError)
		return
	}

	var phonebookData map[string]PhonebookData
	if err := json.Unmarshal(file, &phonebookData); err != nil {
		http.Error(w, "Error parsing phonebook data", http.StatusInternalServerError)
		return
	}

	// Получение данных для выбранного города
	data, ok := phonebookData[city]
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Получение значений из формы поиска
	nameQuery := r.FormValue("name")
	numberQuery := r.FormValue("number")
	addressQuery := r.FormValue("address")

	// Фильтрация результатов по номеру, названию и адресу больницы
	var results []SchoolEntry
	for _, entry := range data.Schools {
		if (numberQuery != "" && !strings.Contains(strings.ToLower(entry.Number), strings.ToLower(numberQuery))) ||
			(nameQuery != "" && !strings.Contains(strings.ToLower(entry.Name), strings.ToLower(nameQuery))) ||
			(addressQuery != "" && !strings.Contains(strings.ToLower(entry.Address), strings.ToLower(addressQuery))) {
			continue
		}

		results = append(results, entry)
	}

	// Отправка отфильтрованных данных в шаблон для результатов поиска
	tmpl.ExecuteTemplate(w, "schools-result.tmpl", PhonebookData{City: city, Schools: results})

}

//Hospitals--------------------------------------------------------------------------

func searchHospitalsPageHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Path[len("/city/hospitals/"):]
	if city == "" {
		http.NotFound(w, r)
		return
	}

	// Отправка данных в шаблон для новой страницы поиска по больницам
	tmpl.ExecuteTemplate(w, "search-hospitals.tmpl", PhonebookData{City: city})
}

func searchHospitalsHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Path[len("/search/hospitals/"):]
	if city == "" {
		http.NotFound(w, r)
		return
	}

	// Чтение данных из JSON
	file, err := ioutil.ReadFile("phonebook.json")
	if err != nil {
		http.Error(w, "Error reading phonebook data", http.StatusInternalServerError)
		return
	}

	var phonebookData map[string]PhonebookData
	if err := json.Unmarshal(file, &phonebookData); err != nil {
		http.Error(w, "Error parsing phonebook data", http.StatusInternalServerError)
		return
	}

	// Получение данных для выбранного города
	data, ok := phonebookData[city]
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Получение значений из формы поиска
	nameQuery := r.FormValue("name")
	numberQuery := r.FormValue("number")
	addressQuery := r.FormValue("address")

	// Фильтрация результатов по номеру, названию и адресу больницы
	var results []HospitalEntry
	for _, entry := range data.Hospitals {
		if (numberQuery != "" && !strings.Contains(strings.ToLower(entry.Number), strings.ToLower(numberQuery))) ||
			(nameQuery != "" && !strings.Contains(strings.ToLower(entry.Name), strings.ToLower(nameQuery))) ||
			(addressQuery != "" && !strings.Contains(strings.ToLower(entry.Address), strings.ToLower(addressQuery))) {
			continue
		}

		results = append(results, entry)
	}

	// Отправка отфильтрованных данных в шаблон для результатов поиска
	tmpl.ExecuteTemplate(w, "hospitals-result.tmpl", PhonebookData{City: city, Hospitals: results})
}

//----------------------------------------------------------------------------------

func searchPageHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Path[len("/city/"):]
	if city == "" {
		http.NotFound(w, r)
		return
	}

	// Отправка данных в шаблон для новой страницы поиска
	tmpl.ExecuteTemplate(w, "search-page.tmpl", PhonebookData{City: city})
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Path[len("/search/"):]
	if city == "" {
		http.NotFound(w, r)
		return
	}

	// Чтение данных из JSON
	file, err := ioutil.ReadFile("phonebook.json")
	if err != nil {
		http.Error(w, "Error reading phonebook data", http.StatusInternalServerError)
		return
	}

	var phonebookData map[string]PhonebookData
	if err := json.Unmarshal(file, &phonebookData); err != nil {
		http.Error(w, "Error parsing phonebook data", http.StatusInternalServerError)
		return
	}

	// Получение данных для выбранного города
	data, ok := phonebookData[city]
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Получение значений из формы поиска
	phoneQuery := r.FormValue("phone")
	nameQuery := r.FormValue("name")

	// Фильтрация результатов по номеру телефона или фамилии, имени, отчеству
	var results []PhoneEntry
	for _, entry := range data.Phonebook {
		// Проверка по номеру телефона
		if phoneQuery != "" && !strings.Contains(entry.Phone, phoneQuery) {
			continue
		}

		// Проверка по фамилии, имени, отчеству (без учета регистра)
		if nameQuery != "" && !strings.Contains(strings.ToLower(entry.Name), strings.ToLower(nameQuery)) {
			continue
		}

		results = append(results, entry)
	}

	// Отправка отфильтрованных данных в шаблон для результатов поиска
	tmpl.ExecuteTemplate(w, "search-result.tmpl", PhonebookData{City: city, Phonebook: results})

}
