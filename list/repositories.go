package list

import (
	"database/sql"
	"log"
)

type List struct {
	Name        string
	Description string
	UserId      int
	Id          int
}

type ListRepository struct {
	db *sql.DB
}

func NewListRepository(db *sql.DB) ListRepository {
	return ListRepository{db: db}
}

func (lr ListRepository) CreateList(model List) (int, error) {
	rows := lr.db.QueryRow(
		"insert into list (name, description, userid) values (?,?, ?) returning id",
		model.Name, model.Description, model.UserId,
	)

	var id int
	err := rows.Scan(&id)
	return id, err
}

func (lr ListRepository) UpdateList(model List) (bool, error) {
	_, err := lr.db.Exec(
		"update list (name, description) values (?,?) where id = ? and userid = ?",
		model.Name, model.Description, model.Id, model.UserId,
	)
	if err != nil {
		log.Println("🤔 [UpdateList] query failed to execute")
		return false, err
	}

	return true, nil
}

func (lr ListRepository) DeleteList(listId, userId int) bool {
	_, err := lr.db.Exec("delete from list where id = ? and userid = ?", listId, userId)
	if err != nil {
		log.Println("🤔 [DeleteList] query failed to execute", err)
		return false
	}
	return true
}

func (lr ListRepository) GetList(listId, userId int) (List, error) {
	row := lr.db.QueryRow("select id, name, description, userid from list where id = id and userid = ?",
		listId, userId)

	var list List

	err := row.Scan(&list.Id, &list.Name, &list.Description, &list.UserId)
	return list, err
}

func (lr ListRepository) GetAllLists(userId int) ([]List, error) {
	rows, err := lr.db.Query("select id, name, description, userid from list where userid = ?", userId)
	if err != nil {
		log.Println("🤔 [GetAllLists] query failed", err)
		return []List{}, nil
	}

	var lists []List

	for rows.Next() {
		var list List
		err := rows.Scan(&list.Id, &list.Name, &list.Description, &list.UserId)
		if err != nil {
			return []List{}, err
		}
		lists = append(lists, list)
	}
	return lists, nil
}
