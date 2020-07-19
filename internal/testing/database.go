package testing

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL driver
)

type SQLDatabase struct {
	db     *sql.DB
	schema string
	table  string
}

func NewSQLDatabase(dsn, table string) (*SQLDatabase, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open test database: %s", err)
	}

	schema := "public"
	if strings.Contains(table, ".") {
		tableArgs := strings.Split(table, ".")
		schema = tableArgs[0]
		table = tableArgs[1]
	}

	return &SQLDatabase{db: db, schema: schema, table: table}, nil
}

func (db *SQLDatabase) LoadFixture(filename string) error {
	rawSql, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not open fixture: %w", err)
	}

	_, err = db.db.Exec(string(rawSql))
	return err
}

func (db *SQLDatabase) Reset() error {
	_, err := db.db.Exec(fmt.Sprintf(`TRUNCATE "%s"."%s";`, db.schema, db.table))
	return err
}
