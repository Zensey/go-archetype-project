package migrator

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/gobuffalo/pop/v5/logging"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func init() {
	SetLogger(defaultLogger)
}

var mrx = regexp.MustCompile(`^(\d+)_([^.]+)(\.[a-z0-9]+)?\.(up|down)\.(sql|fizz)$`)

// NewMigrator returns a new "blank" migrator. It is recommended
// to use something like MigrationBox or FileMigrator. A "blank"
// Migrator should only be used as the basis for a new type of
// migration system.
func NewMigrator(c *gorm.DB) Migrator {
	return Migrator{
		Connection: c,
		Migrations: map[string]Migrations{
			"up":   {},
			"down": {},
		},
	}
}

// Migrator forms the basis of all migrations systems.
// It does the actual heavy lifting of running migrations.
// When building a new migration system, you should embed this
// type into your migrator.
type Migrator struct {
	Connection *gorm.DB
	SchemaPath string
	Migrations map[string]Migrations
}

func (m Migrator) migrationIsCompatible(c *gorm.DB, mi Migration) bool {
	if mi.DBType == "all" || mi.DBType == c.Config.Dialector.Name() {
		return true
	}
	return false
}

const migrationTableName = "schema_migration"

type MigrationTable struct {
	Version string `gorm:"size:14;index:schema_migration_idx_name,unique"`
}

func (MigrationTable) TableName() string {
	return migrationTableName
}

//
//// UpLogOnly insert pending "up" migrations logs only, without applying the patch.
//// It's used when loading the schema dump, instead of the migrations.
func (m Migrator) UpLogOnly() error {
	c := m.Connection
	return m.exec(func() error {
		mtn := migrationTableName

		mfs := m.Migrations["up"]
		sort.Sort(mfs)
		return c.Transaction(func(tx *gorm.DB) error {
			for _, mi := range mfs {
				if !m.migrationIsCompatible(c, mi) {
					continue
				}
				var cnt int64
				err := c.Where("version = ?", mi.Version).Model(&MigrationTable{}).Count(&cnt).Error
				exists := cnt > 0

				//exists, err := c.Where("version = ?", mi.Version).Exists(mtn)
				if err != nil {
					return errors.Wrapf(err, "problem checking for migration version %s", mi.Version)
				}
				if exists {
					continue
				}
				err = tx.Exec(fmt.Sprintf("insert into %s (version) values ('%s')", mtn, mi.Version)).Error
				if err != nil {
					return errors.Wrapf(err, "problem inserting migration version %s", mi.Version)
				}
			}
			return nil
		})
	})
}

// Up runs pending "up" migrations and applies them to the database.
func (m Migrator) Up() error {
	_, err := m.UpTo(0)
	return err
}

// UpTo runs up to step "up" migrations and applies them to the database.
// If step <= 0 all pending migrations are run.
func (m Migrator) UpTo(step int) (applied int, err error) {
	c := m.Connection

	err = m.exec(func() error {
		mtn := migrationTableName
		mfs := m.Migrations["up"]
		mfs.Filter(func(mf Migration) bool {
			return m.migrationIsCompatible(c, mf)
		})
		sort.Sort(mfs)
		for _, mi := range mfs {

			var cnt int64
			err := c.Where("version = ?", mi.Version).Model(&MigrationTable{}).Count(&cnt).Error
			exists := cnt > 0

			//exists, err := c.Where("version = ?", mi.Version).Exists(mtn)
			if err != nil {
				return errors.Wrapf(err, "problem checking for migration version %s", mi.Version)
			}
			if exists {
				continue
			}
			err = c.Transaction(func(tx *gorm.DB) error {
				err := mi.Run(tx)
				if err != nil {
					return err
				}
				err = tx.Exec(fmt.Sprintf("insert into %s (version) values ('%s')", mtn, mi.Version)).Error
				return errors.Wrapf(err, "problem inserting migration version %s", mi.Version)
			})
			if err != nil {
				return err
			}
			log(logging.Info, "> %s", mi.Name)
			applied++
			if step > 0 && applied >= step {
				break
			}
		}
		if applied == 0 {
			log(logging.Info, "Migrations already up to date, nothing to apply")
		} else {
			log(logging.Info, "Successfully applied %d migrations.", applied)
		}
		return nil
	})
	return
}

// Down runs pending "down" migrations and rolls back the
// database by the specified number of steps.
func (m Migrator) Down(step int) error {
	c := m.Connection
	return m.exec(func() error {
		mtn := migrationTableName
		count := int64(0)
		err := c.Table(mtn).Count(&count).Error
		if err != nil {
			return errors.Wrap(err, "migration down: unable count existing migration")
		}
		mfs := m.Migrations["down"]
		mfs.Filter(func(mf Migration) bool {
			return m.migrationIsCompatible(c, mf)
		})
		sort.Sort(sort.Reverse(mfs))
		// skip all ran migration
		if len(mfs) > int(count) {
			mfs = mfs[len(mfs)-int(count):]
		}
		// run only required steps
		if step > 0 && len(mfs) >= step {
			mfs = mfs[:step]
		}
		for _, mi := range mfs {
			var cnt int64
			err := c.Where("version = ?", mi.Version).Model(&MigrationTable{}).Count(&cnt).Error
			exists := cnt > 0

			//exists, err := c.Where("version = ?", mi.Version).Exists(mtn)
			if err != nil {
				return errors.Wrapf(err, "problem checking for migration version %s", mi.Version)
			}
			if !exists {
				return errors.Errorf("migration version %s does not exist", mi.Version)
			}
			err = c.Transaction(func(tx *gorm.DB) error {
				err := mi.Run(tx)
				if err != nil {
					return err
				}
				err = tx.Exec(fmt.Sprintf("delete from %s where version = ?", mtn), mi.Version).Error
				return errors.Wrapf(err, "problem deleting migration version %s", mi.Version)
			})
			if err != nil {
				return err
			}

			log(logging.Info, "< %s", mi.Name)
		}
		return nil
	})
}

// Reset the database by running the down migrations followed by the up migrations.
func (m Migrator) Reset() error {
	err := m.Down(-1)
	if err != nil {
		return err
	}
	return m.Up()
}

// CreateSchemaMigrations sets up a table to track migrations. This is an idempotent
// operation.
func CreateSchemaMigrations(c *gorm.DB) error {
	mtn := migrationTableName

	err := c.Exec(fmt.Sprintf("select * from %s", mtn)).Error
	if err == nil {
		return nil
	}

	return c.Transaction(func(tx *gorm.DB) error {
		if !c.Migrator().HasTable(&MigrationTable{}) {
			return c.Migrator().CreateTable(&MigrationTable{})
		}
		err := c.Migrator().AutoMigrate(&MigrationTable{})
		if err != nil {
			return err
		}

		return nil
	})
}

// CreateSchemaMigrations sets up a table to track migrations. This is an idempotent
// operation.
func (m Migrator) CreateSchemaMigrations() error {
	return CreateSchemaMigrations(m.Connection)
}

// Status prints out the status of applied/pending migrations.
func (m Migrator) Status(out io.Writer) error {
	err := m.CreateSchemaMigrations()
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', tabwriter.TabIndent)
	_, _ = fmt.Fprintln(w, "Version\tName\tStatus\t")
	for _, mf := range m.Migrations["up"] {
		var cnt int64
		err := m.Connection.Where("version = ?", mf.Version).Model(&MigrationTable{}).Count(&cnt).Error
		exists := cnt > 0

		if err != nil {
			return errors.Wrapf(err, "problem with migration")
		}
		state := "Pending"
		if exists {
			state = "Applied"
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t\n", mf.Version, mf.Name, state)
	}
	return w.Flush()
}

// DumpMigrationSchema will generate a file of the current database schema
// based on the value of Migrator.SchemaPath
func (m Migrator) DumpMigrationSchema() error {

	return nil
}

func (m Migrator) exec(fn func() error) error {
	now := time.Now()
	defer func() {
		err := m.DumpMigrationSchema()
		if err != nil {
			log(logging.Warn, "Migrator: unable to dump schema: %v", err)
		}
	}()
	defer printTimer(now)

	err := m.CreateSchemaMigrations()
	if err != nil {
		return errors.Wrap(err, "Migrator: problem creating schema migrations")
	}
	return fn()
}

func printTimer(timerStart time.Time) {
	diff := time.Since(timerStart).Seconds()
	if diff > 60 {
		log(logging.Info, "%.4f minutes", diff/60)
	} else {
		log(logging.Info, "%.4f seconds", diff)
	}
}
