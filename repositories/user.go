package repositories

import (
	"database/sql"
	"errors"
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

func (ur *UserRepository) GetUser(userName string) (User, error) {
	rows, err := ur.db.Query("SELECT id, username, password from users where username = ?", userName)
	if err != nil {
		log.Println("[GetUser] query failed", err)
		return User{}, err
	}

	users := []User{}
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Username, &user.Password)
		users = append(users, user)
	}

	if len(users) > 1 {
		log.Println("More than one user existed with username", userName)
		return User{}, errors.New("did not expect more than one user in get user query")
	} else if len(users) == 0 {
		log.Println("More than one user existed with username", userName)
		return User{}, errors.New("no users were returned for get user query")
	}

	return users[0], nil
}

func (ur *UserRepository) CreateUser(model User) (int, error) {
	row := ur.db.QueryRow("INSERT INTO users (username, password) VALUES (?,?) returning id", model.Username, model.Password)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, nil
	}

	return id, nil
}
