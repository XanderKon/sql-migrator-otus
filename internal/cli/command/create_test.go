package command

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/XanderKon/sql-migrator-otus/internal/cli/config"
	"github.com/XanderKon/sql-migrator-otus/internal/logger"
)

func TestCreate(t *testing.T) {
	// disable logger for testing
	logger := logger.New("DEBUG", io.Discard)

	tests := []struct {
		cfg              *config.MigratorConf
		filenames        []string
		expectedFileName string
		expectedErr      error
	}{
		{
			&config.MigratorConf{
				DSN:  "",
				Dir:  t.TempDir(),
				Type: "sql",
			},
			[]string{
				"TestMigrationSQL",
				"Test.migration.SQL",
			},
			"test_migration_sql.sql",
			nil,
		},
		{
			&config.MigratorConf{
				DSN:  "",
				Dir:  t.TempDir(),
				Type: "go",
			},
			[]string{
				"TestMigrationGO",
				"Test.migration.GO",
			},
			"test_migration_go.go",
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			cmd := &Create{
				Cfg:    tt.cfg,
				Logger: logger,
			}

			for _, f := range tt.filenames {
				time.Sleep(900 * time.Millisecond)
				err := cmd.create(f)
				if tt.expectedErr == nil && err != nil {
					t.Errorf("Unexpected rrror: %v", err)
				}

				files, err := os.ReadDir(tt.cfg.Dir)
				if err != nil {
					t.Fatal(err)
				}

				// check created files
				for _, f := range files {
					if !strings.Contains(f.Name(), tt.expectedFileName) {
						t.Errorf("Error: Expected contains: %v, but received: %v", tt.expectedFileName, f.Name())
					}
				}
			}
		})
	}
}
