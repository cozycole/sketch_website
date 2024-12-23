package models

import (
	"testing"

	"sketchdb.cozycole.net/internal/assert"
	// "sketchdb.cozycole.net/internal/utils"
)

func TestSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	if err := restoreDbScript("../../sql/testdata/test1.sql"); err != nil {
		t.Fatal(err)
	}

	db := newTestDB(t)
	m := SearchModel{db}

	tests := []struct {
		name       string
		query      string
		wantResLen int
		wantType   string
		wantName   string
		wantSlug   string
	}{
		{
			name:       "Find Kyle",
			query:      "kyle",
			wantResLen: 1,
			wantType:   "person",
			wantName:   "Kyle Mooney",
			wantSlug:   "kyle-mooney-1",
		},
		{
			name:       "Find nathanfielder",
			query:      "nathanfielder",
			wantResLen: 1,
			wantType:   "creator",
			wantName:   "nathanfielder",
			wantSlug:   "nathanfielder-1",
		},
		{
			name:       "Find Pumpkin",
			query:      "David Pumpkin",
			wantResLen: 1,
			wantType:   "character",
			wantName:   "David S. Pumpkins",
			wantSlug:   "david-s-pumpkins-1",
		},
	}

	for _, tt := range tests {
		results, err := m.Search(tt.query)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, len(results), tt.wantResLen)
		if len(results) > 0 {
			res := results[0]
			assert.EqualPointer(t, res.Type, &tt.wantType)
			assert.EqualPointer(t, res.Name, &tt.wantName)
			assert.EqualPointer(t, res.Slug, &tt.wantSlug)
		}
	}
}
