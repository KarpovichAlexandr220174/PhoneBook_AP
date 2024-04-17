// authHandler.go

package main

import (
	"PhoneBook_AP/pkg/drivers"
	"PhoneBook_AP/pkg/models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		hashedPassword, err := models.GetHashedPassword(username)
		if err != nil {
			fmt.Printf("Error getting hashed password: %v\n", err)
			http.Error(w, "Invalid login.css credentials", http.StatusUnauthorized)
			return
		}

		fmt.Printf("Hashed Password from DB: %v\n", hashedPassword)

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				fmt.Printf("Invalid password: %v\n", err)
				http.Error(w, "Invalid password", http.StatusUnauthorized)
				return
			}
			fmt.Printf("Error comparing passwords: %v\n", err)
			http.Error(w, "Invalid login.css credentials", http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/application", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "login.tmpl", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		var exists bool
		err := drivers.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
		if err != nil {
			http.Error(w, "Error checking for existing username", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		err = drivers.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
		if err != nil {
			http.Error(w, "Error checking for existing email", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		user := models.User{Username: username, Password: password, Email: email}
		err = models.CreateUser(user)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "register.tmpl", nil)
}
