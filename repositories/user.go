package repositories

import (
	"database/sql"
	"log"
)

type User struct {
	Username string
	Password string
	Id       int
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (ur *UserRepository) UserExists(userName string) (bool, error) {
	rows, err := ur.db.Query("SELECT username FROM users where username = ?", userName)
	if err != nil {
		log.Println("🤔 [UserExists] query failed", err)
		return false, err
	}

	usernames := []string{}
	for rows.Next() {
		var username string
		rows.Scan(&username)
		usernames = append(usernames, username)
	}

	return len(usernames) > 0, nil
}

func (ur *UserRepository) CreateUser(model User) (bool, error) {
	_, err := ur.db.Exec("INSERT INTO users (username, password) VALUES (?,?)", model.Username, model.Password)
	if err != nil {
		log.Println("🤔 [CreateUser] query failed", err)
		return false, err
	}
	return true, nil
}
