//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Profile = newProfileTable("public", "profile", "")

type profileTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnString
	Email     postgres.ColumnString
	FullName  postgres.ColumnString
	Avatar    postgres.ColumnString
	CreatedAt postgres.ColumnTimestamp
	UpdatedAt postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ProfileTable struct {
	profileTable

	EXCLUDED profileTable
}

// AS creates new ProfileTable with assigned alias
func (a ProfileTable) AS(alias string) *ProfileTable {
	return newProfileTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ProfileTable with assigned schema name
func (a ProfileTable) FromSchema(schemaName string) *ProfileTable {
	return newProfileTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ProfileTable with assigned table prefix
func (a ProfileTable) WithPrefix(prefix string) *ProfileTable {
	return newProfileTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ProfileTable with assigned table suffix
func (a ProfileTable) WithSuffix(suffix string) *ProfileTable {
	return newProfileTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newProfileTable(schemaName, tableName, alias string) *ProfileTable {
	return &ProfileTable{
		profileTable: newProfileTableImpl(schemaName, tableName, alias),
		EXCLUDED:     newProfileTableImpl("", "excluded", ""),
	}
}

func newProfileTableImpl(schemaName, tableName, alias string) profileTable {
	var (
		IDColumn        = postgres.StringColumn("id")
		EmailColumn     = postgres.StringColumn("email")
		FullNameColumn  = postgres.StringColumn("full_name")
		AvatarColumn    = postgres.StringColumn("avatar")
		CreatedAtColumn = postgres.TimestampColumn("created_at")
		UpdatedAtColumn = postgres.TimestampColumn("updated_at")
		allColumns      = postgres.ColumnList{IDColumn, EmailColumn, FullNameColumn, AvatarColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns  = postgres.ColumnList{EmailColumn, FullNameColumn, AvatarColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return profileTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		Email:     EmailColumn,
		FullName:  FullNameColumn,
		Avatar:    AvatarColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
