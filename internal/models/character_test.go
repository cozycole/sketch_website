package models

import (
	"testing"

	"sketchdb.cozycole.net/internal/assert"
)

func TestCharacterSearch(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	if err := restoreDbScript("../../sql/testdata/test1.sql"); err != nil {
		t.Fatal(err)
	}
	m := CharacterModel{db}

	tests := []struct {
		name        string
		query       string
		resultCount int
		wantSlug    string
		wantName    string
		wantLast    string
		wantImg     string
	}{
		{
			name:        "Find David",
			query:       "Davi",
			resultCount: 1,
			wantName:    "David S. Pumpkins",
			wantSlug:    "david-s-pumpkins-1",
		},
		{
			name:        "Find Daves",
			query:       "Dav",
			resultCount: 2,
			wantName:    "Dave",
			wantSlug:    "dave-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			characters, err := m.Search(tt.query)
			assert.NilError(t, err)

			if len(characters) == 0 {
				t.Fatalf("No person found with query %s", tt.query)
			}

			assert.Equal(t, len(characters), tt.resultCount)
			c := characters[0]
			assert.EqualPointer(t, c.Slug, &tt.wantSlug)
			assert.EqualPointer(t, c.Name, &tt.wantName)
		})
	}
}
