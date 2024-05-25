package migrations

type Migration interface {
	Run() string
}

var Migrations = []Migration{
	create_users{},
	create_list{},
	create_task{},
}
