package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// SQLiteStorage stores records in a SQLite database.
type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		platform TEXT NOT NULL,
		type TEXT NOT NULL,
		data TEXT NOT NULL,
		collected_at DATETIME NOT NULL
	)`)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) Save(record Record) error {
	data, err := json.Marshal(record.Data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	_, err = s.db.Exec(
		"INSERT INTO records (platform, type, data, collected_at) VALUES (?, ?, ?, ?)",
		record.Platform, record.Type, string(data), record.CollectedAt.Format(time.RFC3339),
	)
	return err
}

func (s *SQLiteStorage) Query(platform, recordType string, limit int) ([]Record, error) {
	query := "SELECT platform, type, data, collected_at FROM records WHERE platform = ? AND type = ? ORDER BY collected_at DESC"
	args := []interface{}{platform, recordType}
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query records: %w", err)
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var r Record
		var dataStr, timeStr string
		if err := rows.Scan(&r.Platform, &r.Type, &dataStr, &timeStr); err != nil {
			continue
		}
		json.Unmarshal([]byte(dataStr), &r.Data)
		r.CollectedAt, _ = time.Parse(time.RFC3339, timeStr)
		records = append(records, r)
	}
	return records, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
