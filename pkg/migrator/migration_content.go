package migrator

import (
	"io"
	"io/ioutil"

	"gorm.io/gorm"
)

// MigrationContent returns the content of a migration.
func MigrationContent(mf Migration, c *gorm.DB, r io.Reader, usingTemplate bool) (string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", nil
	}
	return string(b), nil
}
