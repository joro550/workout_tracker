package migrations

type create_users struct{}

func (cu create_users) Run() string {
	return `
    CREATE TABLE IF NOT EXISTS users (
        Id integer primary key autoincrement,
        Username varchar(250),
        Password varchar(250)
    )
    `
}
