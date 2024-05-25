package users

import (
	"database/sql"
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
	rows := ur.db.QueryRow("SELECT count(id) FROM users where username = ?", userName)

	var count int
	rows.Scan(&count)
	return count > 0, nil
}

func (ur *UserRepository) GetUser(userName string) (User, error) {
	rows := ur.db.QueryRow("SELECT id, username, password from users where username = ?", userName)

	var user User
	err := rows.Scan(&user.Id, &user.Username, &user.Password)
	return user, err
}

func (ur *UserRepository) CreateUser(model User) (int, error) {
	row := ur.db.QueryRow("INSERT INTO users (username, password) VALUES (?,?) returning id", model.Username, model.Password)

	var id int
	err := row.Scan(&id)
	return id, err
}
