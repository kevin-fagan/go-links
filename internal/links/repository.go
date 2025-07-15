package links

import (
	"database/sql"
	"errors"

	"github.com/kevin-fagan/go-links/internal/db"
	"github.com/kevin-fagan/go-links/internal/logs"
)

var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrLinkAlreadyExists = errors.New("link already exists")
)

type Repository struct {
	*db.SQLiteContext
	logs logs.Repository
}

func NewRepository(ctx *db.SQLiteContext) *Repository {
	return &Repository{ctx, *logs.NewRepository(ctx)}
}

func (r *Repository) CountVisit(short string) error {
	statement := `
		UPDATE links
		SET visits = visits + 1
		WHERE short_url = ?;`

	results, err := r.Exec(statement, short)
	if err != nil {
		return err
	}

	rows, _ := results.RowsAffected()
	if rows == 0 {
		return ErrLinkNotFound
	}

	return nil
}

func (r *Repository) Read(short string) (*Link, error) {
	statement := `
		SELECT short_url, long_url, visits, last_updated
		FROM links
		WHERE short_url = ?`

	var link Link

	row := r.QueryRow(statement, short)
	err := row.Scan(&link.ShortURL, &link.LongURL, &link.Visits, &link.LastUpdated)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *Repository) ReadAll(page, pageSize int, search string) ([]Link, int, error) {
	var (
		count int
		links []Link
		err   error
	)

	err = r.WithTx(func(tx *sql.Tx) error {
		links, err = r.ReadAllTx(tx, page, pageSize, search)
		if err != nil {
			return err
		}

		count, err = r.CountTx(tx, search)
		if err != nil {
			return err
		}

		return nil
	})

	return links, count, nil
}

func (r *Repository) Create(short, long, clientIP string) error {
	return r.WithTx(func(tx *sql.Tx) error {
		return r.CreateTx(tx, short, long, clientIP)
	})
}

func (r *Repository) Delete(short, clientIP string) error {
	return r.WithTx(func(tx *sql.Tx) error {
		return r.DeleteTx(tx, short, clientIP)
	})
}

func (r *Repository) Update(short, long, clientIP string) error {
	return r.WithTx(func(tx *sql.Tx) error {
		return r.UpdateTx(tx, short, long, clientIP)
	})
}

func (r *Repository) ReadAllTx(tx *sql.Tx, page, pageSize int, search string) ([]Link, error) {
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

	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ShortURL, &link.LongURL, &link.Visits, &link.LastUpdated)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, nil
}

func (r *Repository) CreateTx(tx *sql.Tx, short, long, clientIP string) error {
	_, err := tx.Exec(`
		INSERT INTO links (short_url, long_url)
		VALUES (?, ?);`, short, long)

	if err != nil {
		return err
	}

	return r.logs.CreateTx(tx, short, long, "", clientIP, "CREATE")
}

func (r *Repository) DeleteTx(tx *sql.Tx, short, clientIP string) error {
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

	return r.logs.CreateTx(tx, short, long, "", clientIP, "DELETE")
}

func (r *Repository) UpdateTx(tx *sql.Tx, short, long, clientIP string) error {
	_, err := tx.Exec(`
		UPDATE links
		SET long_url = ?
		WHERE short_url = ?;`, long, short)

	if err != nil {
		return err
	}

	return r.logs.CreateTx(tx, short, long, "", clientIP, "UPDATE")
}

func (r *Repository) CountTx(tx *sql.Tx, search string) (int, error) {
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
