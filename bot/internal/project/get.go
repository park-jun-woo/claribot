package project

import (
	"database/sql"
	"fmt"

	"parkjunwoo.com/claribot/internal/db"
	"parkjunwoo.com/claribot/internal/types"
)

// Get gets project details
func Get(id string) types.Result {
	globalDB, err := db.OpenGlobal()
	if err != nil {
		return types.Result{
			Success: false,
			Message: fmt.Sprintf("failed to open global db: %v", err),
		}
	}
	defer globalDB.Close()

	var p Project
	err = globalDB.QueryRow(`
		SELECT id, name, path, type, description, status, created_at, updated_at
		FROM projects WHERE id = ?
	`, id).Scan(&p.ID, &p.Name, &p.Path, &p.Type, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		return types.Result{
			Success: false,
			Message: fmt.Sprintf("project not found: %s", id),
		}
	}
	if err != nil {
		return types.Result{
			Success: false,
			Message: fmt.Sprintf("failed to query project: %v", err),
		}
	}

	msg := fmt.Sprintf("Project: %s\nType: %s\nPath: %s\nDescription: %s\nStatus: %s\nCreated: %s",
		p.ID, p.Type, p.Path, p.Description, p.Status, p.CreatedAt)

	return types.Result{
		Success: true,
		Message: msg,
		Data:    &p,
	}
}
