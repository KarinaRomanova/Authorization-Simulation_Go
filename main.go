package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

// UserData представляет данные пользователя
type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Database представляет псевдо БД
type Database struct {
	data map[string]string
	mu   sync.Mutex
}

func main() {

	db := &Database{
		data: map[string]string{
			"user1": "password1",
			"user2": "password2",
		},
	}

	// Обработчик для авторизации
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			username := r.Form.Get("username")
			password := r.Form.Get("password")
			userData := UserData{Username: username, Password: password}
			isAuthenticated := db.checkAuthentication(userData.Username, userData.Password)

			response := struct {
				Authenticated bool `json:"authenticated"`
			}{
				Authenticated: isAuthenticated,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Загрузка HTML-страницы
		tmpl := template.Must(template.New("index").Parse(htmlTemplate))
		tmpl.Execute(w, nil)
	})

	// Обработчик для обновления БД
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			username := r.FormValue("username")
			password := r.FormValue("password")

			// Установим кодировку UTF-8 для ответа
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")

			userData := UserData{Username: username, Password: password}
			db.updateUserData(userData.Username, userData.Password)

			fmt.Fprintf(w, "Данные пользователя %s обновлены", userData.Username)
			return
		}

		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	})

	// Запуск сервера на порту 8080
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func (db *Database) checkAuthentication(username, password string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	storedPassword, exists := db.data[username]
	if !exists || storedPassword != password {
		return false
	}
	return true
}

func (db *Database) updateUserData(username, password string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[username] = password
}

// HTML-шаблон
var htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Простая авторизация</title>
    <meta charset="utf-8">
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
    <style>
        body {
            background-color: #f8f9fa;
        }

        .container {
            margin-top: 100px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="row">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title text-center mb-4">Форма входа</h5>
                        <form action="/" method="post">
                            <div class="form-group">
                                <input type="text" class="form-control" name="username" placeholder="Логин" required>
                            </div>
                            <div class="form-group">
                                <input type="password" class="form-control" name="password" placeholder="Пароль" required>
                            </div>
                            <button type="submit" class="btn btn-primary btn-block">Войти</button>
                        </form>
                    </div>
                </div>
            </div>
            <div class="col-md-6">
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title text-center mb-4">Обновление БД</h5>
                        <form action="/update" method="post">
                            <div class="form-group">
                                <input type="text" class="form-control" name="username" placeholder="Логин" required>
                            </div>
                            <div class="form-group">
                                <input type="password" class="form-control" name="password" placeholder="Пароль" required>
                            </div>
                            <button type="submit" class="btn btn-success btn-block">Обновить БД</button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>
`
