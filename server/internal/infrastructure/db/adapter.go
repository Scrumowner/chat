package db

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"server/internal/modules/auth/models"
)

type SqlAdapter struct {
	sql *sqlx.DB
	sq  squirrel.StatementBuilderType
}

func NewSqlAdapter(sql *sqlx.DB) *SqlAdapter {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &SqlAdapter{
		sql: sql,
		sq:  builder,
	}
}

func (s *SqlAdapter) Insert(entity models.Tabler) error {
	fp := GetFieldsAndPointers(entity)
	raw := s.sq.Insert(entity.TableName()).Columns(fp.Fields...)
	values := make([]interface{}, 0)
	for _, v := range fp.Pointers {
		values = append(values, v)
	}
	query, args, err := raw.Values(values...).ToSql()
	log.Println(query, args)
	if err != nil {
		log.Println("Err when create sql query", err)
		return err
	}
	_, err = s.sql.Exec(query, args...)
	if err != nil {
		log.Println("Err when execute sql query", err)
		return err
	}
	return nil
}

func (s *SqlAdapter) Select(entity models.Tabler, cond Condition, dest interface{}) error {
	fp := GetFieldsAndPointers(entity)
	raw := s.sq.Select(fp.Fields...).From(entity.TableName())
	if cond.Eq != nil {
		for key, val := range cond.Eq {
			raw = raw.Where(squirrel.Eq{key: val})
		}
	}
	if cond.NotEq != nil {
		for key, val := range cond.NotEq {
			raw = raw.Where(squirrel.Eq{key: val})
		}
	}
	query, args, err := raw.ToSql()
	if err != nil {
		return err
	}
	err = s.sql.Select(dest, query, args...)
	if err != nil {
		return err
	}
	return nil
}
