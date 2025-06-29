package repository

import (
	"database/sql"
	"errors"

	"github.com/kevin-fagan/go-links/internal/model"
	"github.com/mattn/go-sqlite3"
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

func (l *LinkRepository) IncrementVisits(short string) error {
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
	var statement string
	var count int
	var err error

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
	statement := `
	INSERT INTO links (short_url, long_url)
	VALUES (?, ?);
	`

	_, err := l.sql.Exec(statement, short, long)
	if err == sqlite3.ErrConstraintUnique {
		return ErrLinkAlreadyExists
	}
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkRepository) DeleteLink(short string) error {
	statement := `
	DELETE FROM links 
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

func (l *LinkRepository) UpdateLink(short, long string) error {
	statement := `
	UPDATE links
	SET long_url = ?
	WHERE short_url = ?;
	`

	results, err := l.sql.Exec(statement, long, short)
	if err != nil {
		return err
	}

	rows, _ := results.RowsAffected()
	if rows == 0 {
		return ErrLinkNotFound
	}

	return nil
}
