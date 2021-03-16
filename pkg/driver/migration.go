package driver

const migrationTableName = "schema_migration"

type MigrationTable struct {
	Version string `gorm:"size:14;index:schema_migration_idx_name,unique"`
}

func (t MigrationTable) TableName() string {
	return migrationTableName
}
