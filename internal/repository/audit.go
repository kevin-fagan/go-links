package repository

import (
	"database/sql"

	"github.com/kevin-fagan/go-links/internal/model"
)

func (l *LinkRepository) GetAudit(page, pageSize int, search string) ([]model.Audit, int, error) {
	var (
		count  int
		audits []model.Audit
		err    error
	)

	err = l.withTx(func(tx *sql.Tx) error {
		audits, err = l.getAuditTx(tx, page, pageSize, search)
		if err != nil {
			return err
		}

		count, err = l.countAuditTx(tx, search)
		if err != nil {
			return err
		}

		return nil
	})

	return audits, count, nil
}

func (l *LinkRepository) getAuditTx(tx *sql.Tx, page, pageSize int, search string) ([]model.Audit, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		rows, err = tx.Query(`
		SELECT short_url, long_url, action, client_ip, timestamp
		FROM audit
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?;`, pageSize, pageSize*page)
	} else {
		pattern := "%" + search + "%"
		rows, err = tx.Query(`
			SELECT short_url, long_url, action, client_ip, timestamp
			FROM audit
			WHERE short_url LIKE ? OR long_url LIKE ?
			ORDER BY timestamp DESC
			LIMIT ? OFFSET ?;`, pattern, pattern, pageSize, pageSize*page)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var audits []model.Audit
	for rows.Next() {
		var audit model.Audit
		err := rows.Scan(&audit.ShortURL, &audit.LongURL, &audit.Action, &audit.ClientIP, &audit.Timestamp)
		if err != nil {
			return nil, err
		}

		audits = append(audits, audit)
	}

	return audits, nil
}

func (l *LinkRepository) createAuditTx(tx *sql.Tx, short, long, clientIP, action string) error {
	_, err := tx.Exec(`
	INSERT INTO audit (short_url, long_url, client_ip, action)
	VALUES (?, ?, ?, ?);`, short, long, clientIP, action)

	return err
}

func (l *LinkRepository) countAuditTx(tx *sql.Tx, search string) (int, error) {
	var count int

	if search == "" {
		err := tx.QueryRow(`SELECT COUNT(*) FROM audit;`).Scan(&count)
		return count, err
	} else {
		pattern := "%" + search + "%"

		err := tx.QueryRow(`
		SELECT COUNT(*) FROM audit
		WHERE short_url LIKE ? OR long_url LIKE ?;`, pattern, pattern).Scan(&count)

		return count, err
	}
}
