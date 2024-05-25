package migrations

type create_list struct{}

func (cu create_list) Run() string {
	return `
    CREATE TABLE IF NOT EXISTS list (
        Id integer primary key autoincrement,
        Name varchar(200),
        Description varchar(200),
        UserId INT
    )
    `
}
