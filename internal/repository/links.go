package repository

import (
	"database/sql"
	"errors"

	"github.com/kevin-fagan/go-links/internal/model"
)

var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrLinkAlreadyExists = errors.New("link already exists")
)

type LinkRepository struct {
	sql *SQLContext
}

func NewLinkRepository(ctx *SQLContext) *LinkRepository {
	return &LinkRepository{
		sql: ctx,
	}
}

func (l *LinkRepository) CountLinkVisit(short string) error {
	statement := `
		UPDATE links
		SET visits = visits + 1
		WHERE short_url = ?;`

	results, err := l.sql.Exec(statement, short)
	if err != nil {
		return err
	}

	rows, _ := results.RowsAffected()
	if rows == 0 {
		return ErrLinkNotFound
	}

	return nil
}

func (l *LinkRepository) GetLink(short string) (*model.Link, error) {
	statement := `
		SELECT short_url, long_url, visits, last_updated
		FROM links
		WHERE short_url = ?`

	var link model.Link

	row := l.sql.QueryRow(statement, short)
	err := row.Scan(&link.ShortURL, &link.LongURL, &link.Visits, &link.LastUpdated)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (l *LinkRepository) GetLinks(search string, page, pageSize int) ([]model.Link, int, error) {
	var (
		count int
		links []model.Link
		err   error
	)

	err = l.withTx(func(tx *sql.Tx) error {
		links, err = l.getLinksTx(tx, page, pageSize, search)
		if err != nil {
			return err
		}

		count, err = l.countLinkTx(tx, search)
		if err != nil {
			return err
		}

		return nil
	})

	return links, count, nil
}

func (l *LinkRepository) CreateLink(short, long, clientIP string) error {
	return l.withTx(func(tx *sql.Tx) error {
		return l.createLinkTx(tx, short, long, clientIP)
	})
}

func (l *LinkRepository) DeleteLink(short, clientIP string) error {
	return l.withTx(func(tx *sql.Tx) error {
		return l.deleteLinkTx(tx, short, clientIP)
	})
}

func (l *LinkRepository) UpdateLink(short, long, clientIP string) error {
	return l.withTx(func(tx *sql.Tx) error {
		return l.updateLinkTx(tx, short, long, clientIP)
	})
}

func (l *LinkRepository) withTx(fn func(tx *sql.Tx) error) error {
	tx, err := l.sql.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (l *LinkRepository) getLinksTx(tx *sql.Tx, page, pageSize int, search string) ([]model.Link, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		rows, err = tx.Query(`
		SELECT short_url, long_url, visits, last_updated
		FROM links
		ORDER BY visits DESC
		LIMIT ? OFFSET ?;`, pageSize, pageSize*page)
	} else {
		pattern := "%" + search + "%"
		rows, err = tx.Query(`
			SELECT short_url, long_url, visits, last_updated
			FROM links
			WHERE short_url LIKE ? OR long_url LIKE ?
			ORDER BY visits DESC
			LIMIT ? OFFSET ?;`, pattern, pattern, pageSize, pageSize*page)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var links []model.Link
	for rows.Next() {
		var link model.Link
		err := rows.Scan(&link.ShortURL, &link.LongURL, &link.Visits, &link.LastUpdated)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, nil
}

func (l *LinkRepository) createLinkTx(tx *sql.Tx, short, long, clientIP string) error {
	_, err := tx.Exec(`
		INSERT INTO links (short_url, long_url)
		VALUES (?, ?);`, short, long)

	if err != nil {
		return err
	}

	return l.createAuditTx(tx, short, long, clientIP, "CREATE")
}

func (l *LinkRepository) deleteLinkTx(tx *sql.Tx, short, clientIP string) error {
	var long string
	err := tx.QueryRow(`
		SELECT long_url 
		FROM links 
		WHERE short_url = ?;`, short).Scan(&long)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM links 
		WHERE short_url = ?;`, short)

	if err != nil {
		return err
	}

	return l.createAuditTx(tx, short, long, clientIP, "DELETE")
}

func (l *LinkRepository) updateLinkTx(tx *sql.Tx, short, long, clientIP string) error {
	_, err := tx.Exec(`
		UPDATE links
		SET long_url = ?
		WHERE short_url = ?;`, long, short)

	if err != nil {
		return err
	}

	return l.createAuditTx(tx, short, long, clientIP, "UPDATE")

}

func (l *LinkRepository) countLinkTx(tx *sql.Tx, search string) (int, error) {
	var count int

	if search == "" {
		err := tx.QueryRow(`SELECT COUNT(*) FROM links;`).Scan(&count)
		return count, err
	} else {
		pattern := "%" + search + "%"

		err := tx.QueryRow(`
		SELECT COUNT(*) FROM links
		WHERE short_url LIKE ? OR long_url LIKE ?;`, pattern, pattern).Scan(&count)

		return count, err
	}
}
