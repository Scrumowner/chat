package migrator

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"server/internal/infrastructure/db"
	"server/internal/modules/auth/models"
)

type Migrator struct {
	db *sqlx.DB
}

func NewMigrator(db *sqlx.DB) *Migrator {
	return &Migrator{
		db: db,
	}
}

func (m *Migrator) Migrate(entity models.Tabler) error {
	fp := db.GetFieldsAndPointers(entity)
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", entity.TableName())
	for i, s := range fp.Fields {
		if i == len(fp.Fields)-1 {
			query += fmt.Sprintf("%s %s );", s, fp.FieldsTypes[s][0])
			break
		}
		query += fmt.Sprintf("%s %s,", s, fp.FieldsTypes[s][0])

	}
	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}
	log.Println("SUCSESS", query)
	return nil
}
