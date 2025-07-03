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

	if search == "" {
		err = l.withTx(func(tx *sql.Tx) error {
			links, err = l.getLinksTx(tx, page, pageSize)
			if err != nil {
				return err
			}

			count, err = l.countLinksTx(tx)
			if err != nil {
				return err
			}

			return nil
		})
	} else {
		err = l.withTx(func(tx *sql.Tx) error {
			links, err = l.getLinksSearchTx(tx, search, page, pageSize)
			if err != nil {
				return err
			}

			count, err = l.countLinksSearchTx(tx, search)
			if err != nil {
				return err
			}

			return nil
		})
	}

	if err != nil {
		return nil, 0, err
	}

	return links, count, nil
}

func (l *LinkRepository) CreateLink(short, long string) error {
	return l.withTx(func(tx *sql.Tx) error {
		return l.createLinkTx(tx, short, long)
	})
}

func (l *LinkRepository) DeleteLink(short string) error {
	return l.withTx(func(tx *sql.Tx) error {
		return l.deleteLinkTx(tx, short)
	})
}

func (l *LinkRepository) UpdateLink(short, long string) error {
	return l.withTx(func(tx *sql.Tx) error {
		return l.updateLinkTx(tx, short, long)
	})
}

func (l *LinkRepository) getLinksTx(tx *sql.Tx, page, pageSize int) ([]model.Link, error) {
	var (
		rows *sql.Rows
		err  error
	)

	rows, err = tx.Query(`
		SELECT short_url, long_url, visits, last_updated
		FROM links
		ORDER BY visits DESC
		LIMIT ? OFFSET ?;`, pageSize, pageSize*page)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return l.iterateLinkRows(rows)
}

func (l *LinkRepository) getLinksSearchTx(tx *sql.Tx, search string, page, pageSize int) ([]model.Link, error) {
	var (
		rows *sql.Rows
		err  error
	)

	pattern := "%" + search + "%"
	rows, err = tx.Query(`
			SELECT short_url, long_url, visits, last_updated
			FROM links
			WHERE short_url LIKE ? OR long_url LIKE ?
			ORDER BY visits DESC
			LIMIT ? OFFSET ?;`, pattern, pattern, pageSize, pageSize*page)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return l.iterateLinkRows(rows)
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

func (l *LinkRepository) countLinksTx(tx *sql.Tx) (int, error) {
	var count int
	err := tx.QueryRow(`SELECT COUNT(*) FROM links;`).Scan(&count)

	return count, err
}

func (l *LinkRepository) countLinksSearchTx(tx *sql.Tx, search string) (int, error) {
	var count int
	pattern := "%" + search + "%"
	err := tx.QueryRow(`
		SELECT COUNT(*) FROM links
		WHERE short_url LIKE ? OR long_url LIKE ?;`, pattern, pattern).Scan(&count)

	return count, err
}

func (l *LinkRepository) createLinkTx(tx *sql.Tx, short, long string) error {
	_, err := tx.Exec(`
		INSERT INTO links (short_url, long_url)
		VALUES (?, ?);`, short, long)

	if err != nil {
		return err
	}

	return l.logAuditTx(tx, short, long, "CREATE")
}

func (l *LinkRepository) deleteLinkTx(tx *sql.Tx, short string) error {
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

	return l.logAuditTx(tx, short, long, "DELETE")
}

func (l *LinkRepository) updateLinkTx(tx *sql.Tx, short, long string) error {
	_, err := tx.Exec(`
		UPDATE links
		SET long_url = ?
		WHERE short_url = ?;`, long, short)

	if err != nil {
		return err
	}

	return l.logAuditTx(tx, short, long, "UPDATE")

}

func (l *LinkRepository) logAuditTx(tx *sql.Tx, short, long, action string) error {
	_, err := tx.Exec(`
	INSERT INTO audit (short_url, long_url, action)
	VALUES (?, ?, ?);`, short, long, action)

	return err
}

func (l *LinkRepository) iterateLinkRows(rows *sql.Rows) ([]model.Link, error) {
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
