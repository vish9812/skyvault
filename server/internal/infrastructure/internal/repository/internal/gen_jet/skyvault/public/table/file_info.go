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

var FileInfo = newFileInfoTable("public", "file_info", "")

type fileInfoTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnString
	OwnerID   postgres.ColumnString
	FolderID  postgres.ColumnString
	Name      postgres.ColumnString
	Size      postgres.ColumnInteger
	Extension postgres.ColumnString
	MimeType  postgres.ColumnString
	Category  postgres.ColumnString
	Preview   postgres.ColumnString
	TrashedAt postgres.ColumnTimestamp
	CreatedAt postgres.ColumnTimestamp
	UpdatedAt postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type FileInfoTable struct {
	fileInfoTable

	EXCLUDED fileInfoTable
}

// AS creates new FileInfoTable with assigned alias
func (a FileInfoTable) AS(alias string) *FileInfoTable {
	return newFileInfoTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FileInfoTable with assigned schema name
func (a FileInfoTable) FromSchema(schemaName string) *FileInfoTable {
	return newFileInfoTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FileInfoTable with assigned table prefix
func (a FileInfoTable) WithPrefix(prefix string) *FileInfoTable {
	return newFileInfoTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FileInfoTable with assigned table suffix
func (a FileInfoTable) WithSuffix(suffix string) *FileInfoTable {
	return newFileInfoTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFileInfoTable(schemaName, tableName, alias string) *FileInfoTable {
	return &FileInfoTable{
		fileInfoTable: newFileInfoTableImpl(schemaName, tableName, alias),
		EXCLUDED:      newFileInfoTableImpl("", "excluded", ""),
	}
}

func newFileInfoTableImpl(schemaName, tableName, alias string) fileInfoTable {
	var (
		IDColumn        = postgres.StringColumn("id")
		OwnerIDColumn   = postgres.StringColumn("owner_id")
		FolderIDColumn  = postgres.StringColumn("folder_id")
		NameColumn      = postgres.StringColumn("name")
		SizeColumn      = postgres.IntegerColumn("size")
		ExtensionColumn = postgres.StringColumn("extension")
		MimeTypeColumn  = postgres.StringColumn("mime_type")
		CategoryColumn  = postgres.StringColumn("category")
		PreviewColumn   = postgres.StringColumn("preview")
		TrashedAtColumn = postgres.TimestampColumn("trashed_at")
		CreatedAtColumn = postgres.TimestampColumn("created_at")
		UpdatedAtColumn = postgres.TimestampColumn("updated_at")
		allColumns      = postgres.ColumnList{IDColumn, OwnerIDColumn, FolderIDColumn, NameColumn, SizeColumn, ExtensionColumn, MimeTypeColumn, CategoryColumn, PreviewColumn, TrashedAtColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns  = postgres.ColumnList{OwnerIDColumn, FolderIDColumn, NameColumn, SizeColumn, ExtensionColumn, MimeTypeColumn, CategoryColumn, PreviewColumn, TrashedAtColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return fileInfoTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		OwnerID:   OwnerIDColumn,
		FolderID:  FolderIDColumn,
		Name:      NameColumn,
		Size:      SizeColumn,
		Extension: ExtensionColumn,
		MimeType:  MimeTypeColumn,
		Category:  CategoryColumn,
		Preview:   PreviewColumn,
		TrashedAt: TrashedAtColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
