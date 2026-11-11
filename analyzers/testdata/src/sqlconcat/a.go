package sqlconcat

import (
	"context"
	"database/sql"
	"fmt"
)

const QUERY = "SELECT 1 WHERE id = ?"

func bad(db *sql.DB, id string) {
	db.Query("SELECT * FROM t WHERE id = " + id)                     // want `SQL query must be a constant`
	db.Exec(fmt.Sprintf("DELETE FROM t WHERE id = %s", id))          // want `SQL query must be a constant`
	db.QueryContext(context.Background(), "SELECT 1 WHERE id = "+id) // want `SQL query must be a constant`
}

func good(db *sql.DB, id string) {
	db.Query("SELECT * FROM t WHERE id = ?", id)
	db.Query(QUERY, id)
}
