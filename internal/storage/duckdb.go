package storage

import (
	"database/sql"
	"fmt"

	"github.com/AMathur20/Home_Network/internal/models"
	_ "github.com/marcboeker/go-duckdb"
)

type DuckDBStorage struct {
	db *sql.DB
}

func NewDuckDBStorage(path string) (*DuckDBStorage, error) {
	db, err := sql.Open("duckdb", path)
	if err != nil {
		return nil, err
	}

	// Initialize schema
	queries := []string{
		`CREATE TABLE IF NOT EXISTS interface_metrics (
			device_name TEXT,
			interface_name TEXT,
			timestamp TIMESTAMP,
			in_octets UBIGINT,
			out_octets UBIGINT,
			in_speed DOUBLE,
			out_speed DOUBLE,
			status TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON interface_metrics (timestamp)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("failed to execute query %s: %w", q, err)
		}
	}

	return &DuckDBStorage{db: db}, nil
}

func (s *DuckDBStorage) SaveMetric(m models.InterfaceMetric) error {
	_, err := s.db.Exec(`
		INSERT INTO interface_metrics (device_name, interface_name, timestamp, in_octets, out_octets, in_speed, out_speed, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		m.DeviceName, m.InterfaceName, m.Timestamp, m.InOctets, m.OutOctets, m.InSpeed, m.OutSpeed, m.Status)
	return err
}

func (s *DuckDBStorage) GetLatestMetrics() ([]models.InterfaceMetric, error) {
	rows, err := s.db.Query(`
		SELECT device_name, interface_name, timestamp, in_octets, out_octets, in_speed, out_speed, status
		FROM interface_metrics
		QUALIFY ROW_NUMBER() OVER(PARTITION BY device_name, interface_name ORDER BY timestamp DESC) = 1`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.InterfaceMetric
	for rows.Next() {
		var m models.InterfaceMetric
		err := rows.Scan(&m.DeviceName, &m.InterfaceName, &m.Timestamp, &m.InOctets, &m.OutOctets, &m.InSpeed, &m.OutSpeed, &m.Status)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (s *DuckDBStorage) GetMetricHistory(deviceName, interfaceName string, limit int) ([]models.InterfaceMetric, error) {
	rows, err := s.db.Query(`
		SELECT device_name, interface_name, timestamp, in_octets, out_octets, in_speed, out_speed, status
		FROM interface_metrics
		WHERE device_name = ? AND interface_name = ?
		ORDER BY timestamp DESC
		LIMIT ?`, deviceName, interfaceName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.InterfaceMetric
	for rows.Next() {
		var m models.InterfaceMetric
		err := rows.Scan(&m.DeviceName, &m.InterfaceName, &m.Timestamp, &m.InOctets, &m.OutOctets, &m.InSpeed, &m.OutSpeed, &m.Status)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (s *DuckDBStorage) Close() error {
	return s.db.Close()
}
