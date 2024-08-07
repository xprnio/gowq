package database

import "database/sql"

type WorkItem struct {
	Id   int
	Name string

	Order       int
	IsCompleted bool

	// private
	isDeleted bool
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
      "ordering",
      "is_completed"
    )
    VALUES (
      @name,
      @order,
      @is_completed
    )
    RETURNING
      "id",
      "name",
      "ordering",
      "is_completed"
  `
	sqlUpdateWorkItem = `
    UPDATE work_items
    SET 
      "name"          = @name,
      "ordering"      = @order,
      "is_completed"  = @is_completed
    WHERE "id" = @id
    RETURNING
      "id",
      "name",
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
				sql.Named("order", item.Order),
				sql.Named("is_completed", item.IsCompleted),
			)
		default:
			_, err = tx.Exec(
				sqlCreateWorkItem,
				sql.Named("name", item.Name),
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
      "is_completed"  INTEGER NOT NULL DEFAULT 0
    );
  `)
	return err
}
