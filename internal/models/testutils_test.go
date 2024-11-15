package models

import (
	"context"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func newTestDB(t *testing.T) *pgxpool.Pool {
	godotenv.Load()

	db, err := pgxpool.New(context.Background(), os.Getenv("TEST_DB_URL"))
	if err != nil {
		t.Fatal(err)
	}

	schemaDirPath := "../../sql/schema"
	schemaExecOrder := []string{
		"person_table.sql",
		"creator_table.sql",
		"character_table.sql",
		"tag_table.sql",
		"video_table.sql",
	}
	for _, scriptName := range schemaExecOrder {
		scriptPath := path.Join(schemaDirPath, scriptName)
		script, err := os.ReadFile(scriptPath)
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(context.Background(), string(script))
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("../../sql/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(context.Background(), string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})
	return db
}

func restoreDbScript(path string) error {
	return exec.Command("psql", os.Getenv("TEST_DB_URL"), "-f", path).Run()
}

func getStringPtr(value string) *string {
	s := value
	return &s
}
