package sqlite_test

import (
	"context"
	"testing"

	"gooji/internal/repository/sqlite"
)

func TestNewDB(t *testing.T) {
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	// Verify tables exist
	ctx := context.Background()
	var tableName string
	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name='videos'").Scan(&tableName)
	if err != nil {
		t.Errorf("videos table not found: %v", err)
	}
}
