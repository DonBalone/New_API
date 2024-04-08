package main

import (
	"New_API/internal/http-server/server"
	"New_API/internal/storage"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	db, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	user, err := storage.NewUser(db, "a@example.com", "a")
	if err != nil {
		if err.Error() == "invalid email" {
			SendError()
		} else {
			log.Fatal(err)
		}

	}

	NewServer(db, user)

}

func NewServer(db *sql.DB, user *storage.User) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		SendEmails(w, r, db)
		if user != nil {
			//SendUser(w, r, db, user.Email)
			SendChangePassword(w, r, db, user.Email)
			//SendDeleteUser(w, r, db, "asdlo@example.com")
		}

	})
	config := server.NewConfig()

	http.ListenAndServe(config.HTTPServer.Address, nil)
}

func SendEmails(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	emails, err := storage.GetEmails(db)
	if err != nil {
		errorData, err := json.Marshal(err.Error())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(errorData)
		return
	}
	if len(emails) == 0 {
		http.Error(w, "Emails not found", http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(emails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем список email на фронтенд
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func SendUser(w http.ResponseWriter, r *http.Request, db *sql.DB, email string) {
	user, err := storage.FindByEmail(db, email)
	if err != nil {
		errorData, err := json.Marshal(err.Error())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(errorData)
		return
	}
	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func SendChangePassword(w http.ResponseWriter, r *http.Request, db *sql.DB, email string) {
	err := storage.ChangePassword(db, email, "a", "123456")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		errorData, err := json.Marshal(err.Error())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(errorData)
		return
	}
	jsonData, err := json.Marshal(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func SendDeleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB, email string) {
	err := storage.DeleteUser(db, email)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		errorData, err := json.Marshal(err.Error())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(errorData)
		return
	}
	jsonData, err := json.Marshal(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func SendError() {
	http.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		errorData, err := json.Marshal("invalid email")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(errorData)
	})
}
