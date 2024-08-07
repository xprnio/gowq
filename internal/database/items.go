package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
)

type WorkItemTags []string
type WorkItem struct {
	Id   int
	Name string
	Tags WorkItemTags

	Order       int
	IsCompleted bool

	// private
	isDeleted bool
}

func (t *WorkItemTags) Scan(value interface{}) error {
	if value, ok := value.(string); ok {
		*t = make([]string, 0)

		values := strings.Split(value, ",")
		for _, tag := range values {
			tag = strings.TrimSpace(tag)
			tag = strings.ToLower(tag)
			if len(tag) > 0 {
				*t = append(*t, tag)
			}
		}

		return nil
	}

	return fmt.Errorf("could not scan %T as WorkItemTags", value)
}

func (t WorkItemTags) Value() (driver.Value, error) {
	return strings.Join(t, ","), nil
}

func (w *WorkItem) Deleted() *WorkItem {
	w.isDeleted = true
	return w
}

func (db *Database) GetAllItems() ([]WorkItem, error) {
	rows, err := db.conn.Query(`
    SELECT
      "id",
      "name",
      "tags",
      "ordering",
      "is_completed"
    FROM work_items
    ORDER BY
      "ordering"    ASC,
      "name"        ASC,
      "created_at"  ASC
  `)

	if err != nil {
		return nil, err
	}

	items := make([]WorkItem, 0)
	for rows.Next() {
		var item WorkItem
		err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Tags,
			&item.Order,
			&item.IsCompleted,
		)

		if err == nil {
			items = append(items, item)
			continue
		}

		return nil, err
	}

	return items, nil
}

const (
	sqlCreateWorkItem = `
    INSERT INTO work_items (
      "name",
      "tags",
      "ordering",
      "is_completed"
    )
    VALUES (
      @name,
      @tags,
      @order,
      @is_completed
    )
    RETURNING
      "id",
      "name",
      "tags",
      "ordering",
      "is_completed"
  `
	sqlUpdateWorkItem = `
    UPDATE work_items
    SET 
      "name"          = @name,
      "tags"          = @tags,
      "ordering"      = @order,
      "is_completed"  = @is_completed
    WHERE "id" = @id
    RETURNING
      "id",
      "name",
      "tags"
      "ordering",
      "is_completed"
  `
	sqlDeleteWorkItem = `
    DELETE FROM work_items
    WHERE "id" = @id
  `
)

func (db *Database) SaveAllItems(items []WorkItem) ([]WorkItem, error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return items, err
	}

	for _, item := range items {
		var err error

		switch {
		case item.isDeleted:
			_, err = tx.Exec(
				sqlDeleteWorkItem,
				sql.Named("id", item.Id),
			)
		case item.Id != 0:
			_, err = tx.Exec(
				sqlUpdateWorkItem,
				sql.Named("id", item.Id),
				sql.Named("name", item.Name),
				sql.Named("tags", item.Tags),
				sql.Named("order", item.Order),
				sql.Named("is_completed", item.IsCompleted),
			)
		default:
			_, err = tx.Exec(
				sqlCreateWorkItem,
				sql.Named("name", item.Name),
				sql.Named("tags", item.Tags),
				sql.Named("order", item.Order),
				sql.Named("is_completed", item.IsCompleted),
			)
		}

		if err != nil {
			defer tx.Rollback()
			return items, err
		}
	}

	if err := tx.Commit(); err != nil {
		return items, err
	}

	return db.GetAllItems()
}

func (db *Database) initWorkItemsTable() error {
	_, err := db.conn.Exec(`
    CREATE TABLE IF NOT EXISTS work_items (
      "id"            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
      "ordering"      INTEGER NOT NULL DEFAULT 0,
      "name"          TEXT    NOT NULL,
      "tags"          TEXT    NOT NULL DEFAULT '',
      "is_completed"  INTEGER NOT NULL DEFAULT 0
    );
  `)
	return err
}
