package migrations

type create_task struct{}

func (cu create_task) Run() string {
	return `
    CREATE TABLE IF NOT EXISTS task (
        Id integer primary key autoincrement,
        Title varchar(200),
        Date datetime,
        Value varchar(250), 

        Type int,
        ListId int,
        UserId int
    )
    `
}
