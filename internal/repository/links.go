package repository

import (
	"errors"

	"github.com/kevin-fagan/learn-gin/internal/model"
)

var (
	ErrLinkNotFound     = errors.New("link not found")
	ErrNegativePage     = errors.New("'page' cannot be negative")
	ErrNegativePageSize = errors.New("'pageSize' cannot be negative")
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
	WHERE short_name = ?; 
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

func (l *LinkRepository) GetLink(short string) (*model.Link, error) {
	statement := `
	SELECT short_name, long_name, visits, last_updated
	FROM links
	WHERE short_name = ?
	`

	var link model.Link
	row := l.sql.QueryRow(statement, short)
	err := row.Scan(&link.ShortURL, &link.LongURL, &link.Visits, &link.LastUpdated)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (l *LinkRepository) GetLinks(page, pageSize int) ([]model.Link, error) {
	statement := `
	SELECT short_name, long_name, visits, last_updated
	FROM links
	LIMIT ? OFFSET ?;
	`

	var links []model.Link

	rows, err := l.sql.Query(statement, pageSize, pageSize*page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	INSERT INTO links (short_name, long_name)
	VALUES (?, ?);
	`

	_, err := l.sql.Exec(statement, short, long)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkRepository) DeleteLink(short string) error {
	statement := `
	DELETE FROM links 
	WHERE short_name = ?;
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
	SET long_name = ?
	WHERE short_name = ?;
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
