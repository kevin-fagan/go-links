package logs

import (
	"database/sql"

	"github.com/kevin-fagan/go-links/internal/db"
)

type Repository struct {
	*db.SQLiteContext
}

func NewRepository(ctx *db.SQLiteContext) *Repository {
	return &Repository{ctx}
}

func (r *Repository) ReadAll(page, pageSize int, search string) ([]Logs, int, error) {
	var (
		count int
		logs  []Logs
		err   error
	)

	err = r.WithTx(func(tx *sql.Tx) error {
		logs, err = r.ReadTx(tx, page, pageSize, search)
		if err != nil {
			return err
		}

		count, err = r.CountTx(tx, search)
		if err != nil {
			return err
		}

		return nil
	})

	return logs, count, nil
}

func (r *Repository) ReadTx(tx *sql.Tx, page, pageSize int, search string) ([]Logs, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		rows, err = tx.Query(`
		SELECT short_url, long_url, tag, action, client_ip, timestamp
		FROM logs
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?;`, pageSize, pageSize*page)
	} else {
		pattern := "%" + search + "%"
		rows, err = tx.Query(`
			SELECT short_url, long_url, tag, action, client_ip, timestamp
			FROM logs
			WHERE short_url LIKE ? OR long_url LIKE ?
			ORDER BY timestamp DESC
			LIMIT ? OFFSET ?;`, pattern, pattern, pageSize, pageSize*page)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var logs []Logs
	for rows.Next() {
		var log Logs
		err := rows.Scan(
			&log.ShortURL,
			&log.LongURL,
			&log.Tag,
			&log.Action,
			&log.ClientIP,
			&log.Timestamp)

		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *Repository) CreateTx(tx *sql.Tx, short, long, tag, clientIP, action string) error {
	_, err := tx.Exec(`
	INSERT INTO logs (short_url, long_url, tag, client_ip, action)
	VALUES (?, ?, ?, ?, ?);`, short, long, tag, clientIP, action)

	return err
}

func (r *Repository) CountTx(tx *sql.Tx, search string) (int, error) {
	var count int

	if search == "" {
		err := tx.QueryRow(`SELECT COUNT(*) FROM logs;`).Scan(&count)
		return count, err
	} else {
		pattern := "%" + search + "%"

		err := tx.QueryRow(`
		SELECT COUNT(*) FROM logs
		WHERE short_url LIKE ? OR long_url LIKE ?;`, pattern, pattern).Scan(&count)

		return count, err
	}
}
