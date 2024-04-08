package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

type User struct {
	ID                int
	Email             string
	Password          string
	encryptedPassword string
}

func New() (*sql.DB, error) {
	config, err := ConfigureFile()
	if err != nil {
		return nil, err
	}
	db, err := Open(config)
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to the database")

	return db, nil
}
func ConfigureFile() (*Config, error) {
	// Чтение файла configs.yaml
	configFile, err := os.ReadFile("Z:/Golang/New_API/configs/configs.yaml")
	if err != nil {
		return nil, err
	}

	// Создание структуры для хранения данных из YAML
	config := NewConfig()
	// Разбор данных из YAML файла
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
func Open(config *Config) (*sql.DB, error) {
	// Формирование строки подключения
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		config.StorageInfo.User, config.StorageInfo.Password, config.StorageInfo.DBName, config.StorageInfo.SSLMode)
	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Проверяем, что у нас есть соединение
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewUser(db *sql.DB, email string, Password string) (*User, error) {
	// validation...
	if !ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email")
	}
	encryptedPassword, err := HashPassword(Password)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("INSERT INTO users (email, encrypted_password) VALUES ($1, $2)", email, encryptedPassword)
	var id int

	err = db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if err != nil {
		return nil, err
	}
	user := User{
		ID:                id,
		Email:             email,
		Password:          Password,
		encryptedPassword: encryptedPassword,
	}
	return &user, nil
}
func ValidateEmail(email string) bool {
	// Шаблон для валидации email
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Проверяем соответствие email шаблону
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}
func GetEmails(db *sql.DB) ([]string, error) {
	var emails []string

	rows, err := db.Query("SELECT email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return emails, nil
}

func ChangePassword(db *sql.DB, email string, Password string, NewPassword string) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil || !exists {
		return fmt.Errorf("user with email %s not found", email)
	}
	var hash string
	err = db.QueryRow("SELECT encrypted_password FROM users WHERE email = $1", email).Scan(&hash)
	if err != nil {
		return err
	}
	// выдает ошибку
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(Password)); err != nil {
		return err
	}
	encryptedPassword, err := HashPassword(NewPassword)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE users SET encrypted_password = $1 WHERE email = $2", encryptedPassword, email)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(db *sql.DB, email string) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil || !exists {
		return fmt.Errorf("user with email %s not found", email)
	}
	_, err = db.Exec("DELETE FROM users WHERE email = $1", email)
	if err != nil {
		return err
	}
	return nil
}

func FindByEmail(db *sql.DB, email string) (*User, error) {
	var user User

	err := db.QueryRow("SELECT id, encrypted_password FROM users WHERE email = $1", email).Scan(&user.ID, &user.encryptedPassword)
	if err != nil {
		return nil, err
	}
	user.Email = email
	return &user, nil
}

func HashPassword(password string) (string, error) {
	// Генерация хеша пароля с использованием bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
