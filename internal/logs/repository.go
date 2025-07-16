package logs

import (
	"database/sql"
	"time"

	"github.com/kevin-fagan/go-links/internal/db"
)

type Logs struct {
	// ShortURL is only populated in the logs when a change is made to a link
	ShortURL string
	// LongURL is only populated in the logs when a change is made to a link
	LongURL string
	// Tag is only populated in the logs when a change is made to a tag
	Tag string
	// Actions describes the operation performed (create, update, delete) on a link or tag
	Action string
	// ClientIP is the IP address of the user who triggered the action. (OAuth not yet supported.)
	ClientIP string
	// Timestamp records when this log entry was created.
	Timestamp time.Time
}

type Repository struct {
	*db.SQLiteContext
}

func NewRepository(ctx *db.SQLiteContext) *Repository {
	return &Repository{ctx}
}

// ReadAll retrieves a set of logs from the repository along with the total matching count.
// The results are paginated based on the provided page number, page size, and optional search query.
// If an error occurs, any changes are rolled back and the error is returned.
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

// ReadAllTx is a SQL transaction that retrieves a set of logs
// The results are paginated based on the provided page number, page size, and optional search query.
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

// CreateTx is a SQL transaction that creates a log entry
func (r *Repository) CreateTx(tx *sql.Tx, short, long, tag, clientIP, action string) error {
	_, err := tx.Exec(`
	INSERT INTO logs (short_url, long_url, tag, client_ip, action)
	VALUES (?, ?, ?, ?, ?);`, short, long, tag, clientIP, action)

	return err
}

// CountTx is a SQL transaction that returns the numbers of results found
// If search is not empty, it will be used as part of the SQL query
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
