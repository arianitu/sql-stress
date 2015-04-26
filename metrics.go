package main

import (
	"database/sql"
)

type TableSize interface {
	GetTableSize() int64
	GetIndexSize() int64
	GetAvgRowSize() int64
	GetRows() int64
}

type MySQLTableSize struct {
	Table  string
	Db     *sql.DB
	status MySQLTableStatus
}

type MySQLTableStatus struct {
	Name          string
	Engine        string
	Version       int64
	RowFormat     string
	Rows          int64
	AvgRowLength  int64
	DataLength    int64
	MaxDataLength int64
	IndexLength   int64
	DataFree      uint64
	AutoIncrement sql.NullInt64
	CreateTime    sql.NullString
	UpdateTime    sql.NullString
	CheckTime     sql.NullString
	Collation     sql.NullString
	Checksum      sql.NullString
	CreateOptions string
	Comment       string
}

func (m *MySQLTableSize) Init() error {
	_, err := m.Db.Exec("ANALYZE TABLE " + m.Table)
	if err != nil {
		return err
	}

	err = m.Db.QueryRow("SHOW TABLE STATUS LIKE '"+m.Table+"'").Scan(
		&m.status.Name,
		&m.status.Engine,
		&m.status.Version,
		&m.status.RowFormat,
		&m.status.Rows,
		&m.status.AvgRowLength,
		&m.status.DataLength,
		&m.status.MaxDataLength,
		&m.status.IndexLength,
		&m.status.DataFree,
		&m.status.AutoIncrement,
		&m.status.CreateTime,
		&m.status.UpdateTime,
		&m.status.CheckTime,
		&m.status.Collation,
		&m.status.Checksum,
		&m.status.CreateOptions,
		&m.status.Comment)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQLTableSize) GetTableSize() int64 {
	return m.status.DataLength
}

func (m *MySQLTableSize) GetIndexSize() int64 {
	return m.status.IndexLength
}

func (m *MySQLTableSize) GetAvgRowSize() int64 {
	return m.status.AvgRowLength
}

func (m *MySQLTableSize) GetRows() int64 {
	return m.status.Rows
}
