package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Postgres struct {
	db     *sql.DB
	schema string
	table  string
}

func NewPostgres(dsn string, table string) (*Postgres, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open postgres storage: %w", err)
	}

	schema := "public"
	if strings.Contains(table, ".") {
		tableArgs := strings.Split(table, ".")
		schema = tableArgs[0]
		table = tableArgs[1]
	}

	return &Postgres{
		db:     db,
		schema: schema,
		table:  table,
	}, nil
}

func (p *Postgres) Bury(funcs []string) error {
	ctx := context.Background()
	c, err := p.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("could not get session: %w", err)
	}

	newFuncs := make(map[string]struct{}, len(funcs))
	for _, v := range funcs {
		newFuncs[v] = struct{}{}
	}

	_, err = c.ExecContext(ctx, fmt.Sprintf(`BEGIN WORK; LOCK TABLE "%s"."%s" IN ACCESS EXCLUSIVE MODE;`, p.schema, p.table))
	if err != nil {
		return fmt.Errorf("could not get lock: %w", err)
	}

	rows, err := c.QueryContext(ctx, fmt.Sprintf(`SELECT function, first_seen_at FROM "%s"."%s";`, p.schema, p.table))
	if err != nil {
		return fmt.Errorf("could not query table: %w", err)
	}

	// Map all functions
	seen := make(map[string]struct{})
	notSeen := make(map[string]struct{})
	for rows.Next() {
		var function string
		var seenAt sql.NullTime
		err := rows.Scan(&function, &seenAt)
		if err != nil {
			return fmt.Errorf("could not scan row: %w", err)
		}

		if seenAt.Valid {
			seen[function] = struct{}{}
		} else {
			notSeen[function] = struct{}{}
		}
	}

	updateArgs := make([]interface{}, 0)
	insertArgs := make([]interface{}, 0)
	for _, f := range funcs {
		_, wasSeen := seen[f]
		_, wasNotSeen := notSeen[f]

		if !wasSeen && !wasNotSeen {
			insertArgs = append(insertArgs, f)
		} else if wasNotSeen {
			updateArgs = append(updateArgs, f)
		}
	}
	updateArgs = append(updateArgs, insertArgs...)

	if len(insertArgs) > 0 {
		var insertPositions []string
		for i := range insertArgs {
			insertPositions = append(insertPositions, fmt.Sprintf("($%d)", i+1))
		}

		placeholders := strings.Join(insertPositions, ",")
		query := fmt.Sprintf(`INSERT INTO "%s"."%s" (function) VALUES %s ON CONFLICT DO NOTHING;`, p.schema, p.table, placeholders)
		_, err := c.ExecContext(ctx, query, insertArgs...)
		if err != nil {
			return fmt.Errorf("could not insert data: %w", err)
		}
	}

	if len(updateArgs) > 0 {
		updatePositions := make([]string, len(updateArgs))
		for i := range updateArgs {
			updatePositions[i] = fmt.Sprintf("$%d", i+2)
		}
		updateArgs = append([]interface{}{time.Now()}, updateArgs...)

		// TODO Chunked queries, postgres had limit of 10000 args in WHERE IN
		placeholders := strings.Join(updatePositions, ",")
		query := fmt.Sprintf(`UPDATE "%s"."%s" SET first_seen_at = $1 WHERE function NOT IN (%s) AND first_seen_at IS NULL;`, p.schema, p.table, placeholders)
		_, err = c.ExecContext(ctx, query, updateArgs...)
		if err != nil {
			return fmt.Errorf("error when updating table: %w", err)
		}
	}

	_, err = c.ExecContext(ctx, `COMMIT WORK;`)
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (p *Postgres) Dig() ([]string, error) {
	rows, err := p.db.Query(fmt.Sprintf(`SELECT function FROM "%s"."%s" WHERE first_seen_at IS NULL`, p.schema, p.table))
	if err != nil {
		return nil, fmt.Errorf("could not query rows: %w", err)
	}
	var funcs []string
	for rows.Next() {
		var function string
		err := rows.Scan(&function)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		funcs = append(funcs, function)
	}

	return funcs, nil
}
