package infra

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DefaultDbName = "imageboard"
)

func GetDbConn(DbHost string, DbPort string, DbName string, DbUser string, DbPassword string) (*gorm.DB, error) {
	connectionString := GetConnString("mysql", DbHost, DbPort, DbName, DbUser, DbPassword)
	dialector := GetDialector("mysql", connectionString)
	conn, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Database connection failed %s", err)
	}

	var tableList []string
	conn.Raw("SHOW TABLES LIKE 'posts'").Scan(&tableList)
	if len(tableList) == 0 {
		CreateDbSchema(conn)
	}

	return conn, nil
}

func GetDialector(db string, connString string) gorm.Dialector {
	if db == "postgresql" {
		return postgres.Open(connString)
	}

	return mysql.Open(connString)
}

func GetConnString(db string, DbHost string, DbPort string, DbName string, DbUser string, DbPassword string) string {
	if db == "postgresql" {
		return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
			DbHost, DbPort, DbName, DbUser, DbPassword)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		DbUser, DbPassword, DbHost, DbPort, DbName)
}

func CreateDbSchema(db *gorm.DB) {
	data, _ := os.ReadFile("./db_schema.my.sql")

	_ = db.Exec(string(data))
}

func DeleteDb(db *gorm.DB) {
	db.Exec("DROP DATABASE " + DefaultDbName)
}

func ClearDb(db *gorm.DB) {
	tables := []string{
		//"pair_to_role",
		//"action_to_object",
		"links",

		"posts",
		"threads",
		"boards",
		"topics",
		//"users",

		//"roles",
		"log_actions",
		//"objects",
		//"actions",
		"media",
	}

	for _, name := range tables {
		db.Exec("DELETE FROM " + name)
	}
}
