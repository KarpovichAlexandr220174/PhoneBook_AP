package main

import (
	"PhoneBook_AP/pkg/models"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
)

func generateRandomPassword() string {
	const (
		upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowerChars   = "abcdefghijklmnopqrstuvwxyz"
		digitChars   = "0123456789"
		specialChars = "!@#$%^&*()-_=+,.?/:;{}[]`~"
		passwordLen  = 12
	)

	chars := upperChars + lowerChars + digitChars + specialChars
	var password string

	for i := 0; i < passwordLen; i++ {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		password += string(chars[randIndex.Int64()])
	}

	return password
}

func forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Получаем почту пользователя из формы
		email := r.FormValue("email")

		// Генерируем новый случайный пароль
		newPassword := generateRandomPassword()

		// Отправляем новый пароль на почту пользователя
		subject := "Восстановление пароля"
		body := fmt.Sprintf("Ваш новый пароль: %s", newPassword)
		err := sendMail(email, subject, body)
		if err != nil {
			http.Error(w, "Ошибка отправки письма", http.StatusInternalServerError)
			return
		}

		// Обновляем пароль пользователя в базе данных
		if err := models.UpdatePasswordByEmail(email, newPassword); err != nil {
			http.Error(w, "Ошибка обновления пароля в базе данных", http.StatusInternalServerError)
			return
		}

		// Показываем пользователю сообщение о том, что письмо отправлено
		http.Redirect(w, r, "/message", http.StatusSeeOther)
	}

	// Если метод запроса не POST, показываем форму для восстановления пароля
	tmpl.ExecuteTemplate(w, "forgot-password.tmpl", nil)
}
