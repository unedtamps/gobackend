package utils

import "fmt"

func LoadPostgresConnString(port int, user, host, pass, db_name string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pass, host, port, db_name,
	)
}

func LoadMySQLConnString(port int, user, host, password, db_name string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		user, password, host, port, db_name)
}

func LoadMySQLConnStringWithDriver(port int, user, host, password, db_name string) string {
	return fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?parseTime=true",
		user, password, host, port, db_name)
}

func LoadSqliteConnString(sql_path string) string {
	return fmt.Sprintf(
		"file:%s?_pragma=foreign_keys(ON)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)",
		sql_path,
	)
}
