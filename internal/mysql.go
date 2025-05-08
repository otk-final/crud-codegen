package internal

import (
	"context"
	"crud-codegen/schema"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const schemaDatasource = "information_schema"

func mysqlQuery(url string, datasourceName string, tableName string) (*schema.Table, error) {

	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	//table
	tableQuery := fmt.Sprintf("select TABLE_COMMENT from %s.TABLES where TABLE_SCHEMA = ? and TABLE_NAME = ?", schemaDatasource)
	row := db.QueryRowContext(ctx, tableQuery, datasourceName, tableName)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var tableComment string
	err = row.Scan(&tableComment)
	if err != nil {
		return nil, err
	}

	//column
	columnQuery := fmt.Sprintf("SELECT COLUMN_NAME ,DATA_TYPE , COLUMN_COMMENT , COLUMN_KEY  FROM %s.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION", schemaDatasource)
	rows, err := db.QueryContext(ctx, columnQuery, datasourceName, tableName)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	columns := make([]*schema.Column, 0)
	for rows.Next() {
		var name, typename, comment, key string
		err := rows.Scan(&name, &typename, &comment, &key)
		if err != nil {
			return nil, err
		}
		columns = append(columns, &schema.Column{
			PK:        key == "PRI",
			Name:      name,
			Type:      typename,
			Comment:   comment,
			Attribute: map[string]string{},
			Enums:     make([]schema.Enum, 0),
		})
	}

	return &schema.Table{
		Name:    tableName,
		Comment: tableComment,
		Columns: columns,
	}, nil
}
