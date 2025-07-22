package tags

import (
	"database/sql"
	"time"

	"github.com/kevin-fagan/go-links/internal/db"
	"github.com/kevin-fagan/go-links/internal/logs"
)

type Tag struct {
	// Name is the label assigned to the tag (e.g., Social).
	Name string
	// LastUpdated indicates the last time the tag was modified
	LastUpdated time.Time
}

type Repository struct {
	db   *db.SQLiteContext
	logs logs.Repository
}

func NewRepository(ctx *db.SQLiteContext) *Repository {
	return &Repository{ctx, *logs.NewRepository(ctx)}
}

// Read will read a single tag from the repository. An error is returned if one occurs
func (r *Repository) Read(name string) (*Tag, error) {
	statement := `
		SELECT name, last_updated
		FROM tags
		WHERE name = ?`

	var tag Tag

	row := r.db.QueryRow(statement, name)
	err := row.Scan(&tag.Name, &tag.LastUpdated)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

// ReadAll retrieves a set of tags from the repository along with the total matching count.
// The results are paginated based on the provided page number, page size, and optional search query.
// If an error occurs, any changes are rolled back and the error is returned.
func (r *Repository) ReadAll(page, pageSize int, search string) ([]Tag, int, error) {
	var (
		count int
		tags  []Tag
		err   error
	)

	err = r.db.WithTx(func(tx *sql.Tx) error {
		tags, err = r.ReadAllTx(tx, page, pageSize, search)
		if err != nil {
			return err
		}

		count, err = r.ResultsTx(tx, search)
		if err != nil {
			return err
		}

		return nil
	})

	return tags, count, nil
}

// Create will create a tag. Addtionally, a log entry will be created reflecting the operation.
// If an error occurs, any changes are rolled back and the error is returned
func (r *Repository) Create(name, clientIP string) error {
	return r.db.WithTx(func(tx *sql.Tx) error {
		return r.CreateTx(tx, name, clientIP)
	})
}

// Delete will delete a tag. Additionally, a log entry will be created reflecting the operation.
// If an error occurs, any changes are rolled back and the error is returned
func (r *Repository) Delete(name, clientIP string) error {
	return r.db.WithTx(func(tx *sql.Tx) error {
		return r.DeleteTx(tx, name, clientIP)
	})
}

// Update will update a tag. Additionally, a log entry will be created reflecting the operation.
// If an error occurs, any changes are rolled back and the error is returned
func (r *Repository) Update(old, new, clientIP string) error {
	return r.db.WithTx(func(tx *sql.Tx) error {
		return r.UpdateTx(tx, old, new, clientIP)
	})
}

// ReadAllTx is a SQL transaction that retrieves a set of tags
// The results are paginated based on the provided page number, page size, and optional search query.
func (r *Repository) ReadAllTx(tx *sql.Tx, page, pageSize int, search string) ([]Tag, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		rows, err = tx.Query(`
		SELECT name, last_updated
		FROM tags
		ORDER BY last_updated DESC
		LIMIT ? OFFSET ?;`, pageSize, pageSize*page)
	} else {
		pattern := "%" + search + "%"
		rows, err = tx.Query(`
			SELECT name, last_updated
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
		err := rows.Scan(&tag.Name, &tag.LastUpdated)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

// CreateTx is a SQL transaction that creates a tag. Additionally, a log entry will
// be created relfecting the operation
func (r *Repository) CreateTx(tx *sql.Tx, name, clientIP string) error {
	_, err := tx.Exec(`
		INSERT INTO tags (name)
		VALUES (?);`, name)

	if err != nil {
		return err
	}

	return r.logs.CreateTx(tx, "", "", name, clientIP, "CREATE")
}

// DeleteTx is a SQL transaction that deletes a tag. Additionally, a log entry will
// be created relfecting the operation
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

// UpdateTx is a SQL transaction that updates a tag. Additionally, a log entry will
// be created relfecting the operation
func (r *Repository) UpdateTx(tx *sql.Tx, old, new, clientIP string) error {
	_, err := tx.Exec(`
		UPDATE tags
		SET name = ?
		WHERE name = ?;`, new, old)

	if err != nil {
		return err
	}

	return r.logs.CreateTx(tx, "", "", new, clientIP, "UPDATE")
}

// ResultsTx is a SQL transaction that returns the numbers of results found.
// If search is not empty, it will be used as part of the SQL query
func (r *Repository) ResultsTx(tx *sql.Tx, search string) (int, error) {
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
