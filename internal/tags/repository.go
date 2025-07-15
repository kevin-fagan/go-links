package tags

import (
	"database/sql"

	"github.com/kevin-fagan/go-links/internal/db"
	"github.com/kevin-fagan/go-links/internal/logs"
)

type Repository struct {
	*db.SQLiteContext
	logs logs.Repository
}

func NewRepository(ctx *db.SQLiteContext) *Repository {
	return &Repository{ctx, *logs.NewRepository(ctx)}
}

func (r *Repository) Read(name string) (*Tag, error) {
	statement := `
		SELECT name, color, last_updated
		FROM tags
		WHERE name = ?`

	var tag Tag

	row := r.QueryRow(statement, name)
	err := row.Scan(&tag.Name, &tag.Color, &tag.LastUpdated)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *Repository) ReadAll(page, pageSize int, search string) ([]Tag, int, error) {
	var (
		count int
		tags  []Tag
		err   error
	)

	err = r.WithTx(func(tx *sql.Tx) error {
		tags, err = r.readTagsTx(tx, page, pageSize, search)
		if err != nil {
			return err
		}

		count, err = r.CountTx(tx, search)
		if err != nil {
			return err
		}

		return nil
	})

	return tags, count, nil
}

func (r *Repository) Create(name, color, clientIP string) error {
	return r.WithTx(func(tx *sql.Tx) error {
		return r.CreateTx(tx, name, color, clientIP)
	})
}

func (r *Repository) Delete(name, color, clientIP string) error {
	return r.WithTx(func(tx *sql.Tx) error {
		return r.DeleteTx(tx, name, clientIP)
	})
}

func (r *Repository) Update(name, color, clientIP string) error {
	return r.WithTx(func(tx *sql.Tx) error {
		return r.UpdateTx(tx, name, color, clientIP)
	})
}

func (r *Repository) readTagsTx(tx *sql.Tx, page, pageSize int, search string) ([]Tag, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		rows, err = tx.Query(`
		SELECT name, color, last_updated
		FROM tags
		ORDER BY last_updated DESC
		LIMIT ? OFFSET ?;`, pageSize, pageSize*page)
	} else {
		pattern := "%" + search + "%"
		rows, err = tx.Query(`
			SELECT name, color, last_updated
			FROM tags
			WHERE name LIKE ? 
			ORDER BY last_updated DESC
			LIMIT ? OFFSET ?;`, pattern, pattern, pageSize, pageSize*page)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.Name, &tag.Color, &tag.LastUpdated)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *Repository) CreateTx(tx *sql.Tx, name, color, clientIP string) error {
	_, err := tx.Exec(`
		INSERT INTO tags (name, color)
		VALUES (?, ?);`, name, color)

	if err != nil {
		return err
	}

	return r.logs.CreateTx(tx, "", "", name, clientIP, "CREATE")
}

func (r *Repository) DeleteTx(tx *sql.Tx, name, clientIP string) error {
	var tagName string
	err := tx.QueryRow(`
		SELECT name 
		FROM tags 
		WHERE name = ?;`, name).Scan(&tagName)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM tags 
		WHERE name = ?;`, tagName)

	if err != nil {
		return err
	}

	return r.logs.CreateTx(tx, "", "", tagName, clientIP, "DELETE")
}

func (r *Repository) UpdateTx(tx *sql.Tx, name, color, clientIP string) error {
	_, err := tx.Exec(`
		UPDATE tags
		SET name = ?
		WHERE color = ?;`, name, color)

	if err != nil {
		return err
	}

	return r.logs.CreateTx(tx, "", "", name, clientIP, "UPDATE")
}

func (r *Repository) CountTx(tx *sql.Tx, search string) (int, error) {
	var count int

	if search == "" {
		err := tx.QueryRow(`SELECT COUNT(*) FROM tags;`).Scan(&count)
		return count, err
	} else {
		pattern := "%" + search + "%"

		err := tx.QueryRow(`
		SELECT COUNT(*) FROM links
		WHERE name LIKE ?;`, pattern, pattern).Scan(&count)

		return count, err
	}
}
