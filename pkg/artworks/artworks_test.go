package artworks

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/VTGare/boe-tea-backend/internal/database"
	"github.com/VTGare/boe-tea-backend/pkg/artworks/options"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	service Service
)

func up() {
	service.InsertOne(context.Background(), &ArtworkInsert{
		Title:  "Boe Tea",
		Author: "vt",
		URL:    "https://github.com/VTGare/boe-tea-go",
		Images: []string{"https://github.com/VTGare/boe-tea-go"},
	})
}

func TestMain(t *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "4.4", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	var db *database.DB
	if err = pool.Retry(func() error {
		db, err = database.Connect("mongodb://localhost:"+resource.GetPort("tcp/27017"), "boe-tea")
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	service = NewService(db, zap.NewExample().Sugar())

	code := t.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestFindOne(t *testing.T) {
	tests := []struct {
		name    string
		filter  *options.FilterOne
		expect  func(*Artwork) bool
		wantErr bool
	}{
		{
			name: "Test 1. ID 1. Success.",
			filter: &options.FilterOne{
				ID: 1,
			},
			expect: func(a *Artwork) bool {
				return a.ID == 1
			},
		},
	}

	for _, test := range tests {
		artwork, err := service.FindOne(context.Background(), test.filter)
		if test.wantErr {
			assert.Error(t, err, test.name)
		}

		assert.True(t, test.expect(artwork))
	}
}
