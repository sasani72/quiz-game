package migrator

import (
	"database/sql"
	"fmt"
	"quiz-game/repository/mysql"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	dbConfig   mysql.Config
	migrations *migrate.FileMigrationSource
}

// TODO - set migration table name
// TODO - add limit to up and down methods

func New(dbConfig mysql.Config) Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}

	return Migrator{dialect: "mysql", dbConfig: dbConfig, migrations: migrations}
}

func (m Migrator) Up() {
	db, err := sql.Open(m.dialect, fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		m.dbConfig.Username, m.dbConfig.Password, m.dbConfig.Host, m.dbConfig.Port, m.dbConfig.DBName))
	if err != nil {
		panic(fmt.Errorf("mysql open failed: %v", err))
	}
	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("mysql migration failed: %v", err))
	}
	fmt.Printf("Applied %d migrations!\n", n)
}

func (m Migrator) Down() {
	db, err := sql.Open(m.dialect, fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		m.dbConfig.Username, m.dbConfig.Password, m.dbConfig.Host, m.dbConfig.Port, m.dbConfig.DBName))
	if err != nil {
		panic(fmt.Errorf("mysql open failed: %v", err))
	}
	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		panic(fmt.Errorf("mysql migration rollback failed: %v", err))
	}
	fmt.Printf("Rollbacked %d migrations!\n", n)
}

func (m Migrator) Status() {
	// add status
}
