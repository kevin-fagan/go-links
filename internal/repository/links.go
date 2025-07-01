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

func (l *LinkRepository) IncVisits(short string) error {
	statement := `
	UPDATE links
	SET visits = visits + 1
	WHERE short_url = ?; 
	`

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

func (l *LinkRepository) Count(search string) (int, error) {
	var (
		count     int
		statement string
		err       error
	)

	if search == "" {
		statement = `SELECT COUNT(*) FROM links;`
		err = l.sql.QueryRow(statement).Scan(&count)
	} else {
		pattern := "%" + search + "%"
		statement = `
			SELECT COUNT(*) FROM links
			WHERE short_url LIKE ? OR long_url LIKE ?`
		err = l.sql.QueryRow(statement, pattern, pattern).Scan(&count)
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (l *LinkRepository) GetLink(short string) (*model.Link, error) {
	statement := `
	SELECT short_url, long_url, visits, last_updated
	FROM links
	WHERE short_url = ?
	`

	var link model.Link
	row := l.sql.QueryRow(statement, short)
	err := row.Scan(&link.ShortURL, &link.LongURL, &link.Visits, &link.LastUpdated)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (l *LinkRepository) GetLinks(search string, page, pageSize int) ([]model.Link, error) {
	var statement string
	var rows *sql.Rows
	var err error

	if search == "" {
		statement = `
			SELECT short_url, long_url, visits, last_updated
			FROM links
			ORDER BY visits DESC
			LIMIT ? OFFSET ?;`
		rows, err = l.sql.Query(statement, pageSize, pageSize*page)
	} else {
		pattern := "%" + search + "%"
		statement = `
			SELECT short_url, long_url, visits, last_updated
			FROM links
			WHERE short_url LIKE ? OR long_url LIKE ?
			ORDER BY visits DESC
			LIMIT ? OFFSET ?;`
		rows, err = l.sql.Query(statement, pattern, pattern, pageSize, pageSize*page)
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

func (l *LinkRepository) CreateLink(short, long string) error {
	tx, err := l.sql.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = l.createLink(tx, short, long)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (l *LinkRepository) DeleteLink(short string) error {
	tx, err := l.sql.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = l.deleteLink(tx, short)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (l *LinkRepository) UpdateLink(short, long string) error {
	tx, err := l.sql.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = l.updateLink(tx, short, long)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (l *LinkRepository) createLink(tx *sql.Tx, short, long string) error {
	_, err := tx.Exec(`
		INSERT INTO links (short_url, long_url)
		VALUES (?, ?);`, short, long)

	if err != nil {
		return err
	}

	return l.logAudit(tx, short, long, model.Action("CREATE"))
}

func (l *LinkRepository) deleteLink(tx *sql.Tx, short string) error {
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

	return l.logAudit(tx, short, long, model.Action("DELETE"))
}

func (l *LinkRepository) updateLink(tx *sql.Tx, short, long string) error {
	_, err := tx.Exec(`
		UPDATE links
		SET long_url = ?
		WHERE short_url = ?;`, long, short)

	if err != nil {
		return err
	}

	return l.logAudit(tx, short, long, model.Action("UPDATE"))

}

func (l *LinkRepository) logAudit(tx *sql.Tx, short, long string, action model.Action) error {
	_, err := tx.Exec(`
	INSERT INTO audit (short_url, long_url, action)
	VALUES (?, ?, ?);`, short, long, action)

	return err
}
