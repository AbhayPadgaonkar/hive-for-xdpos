package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	db.SetMaxOpenConns(1) // SQLite doesn't support concurrent writes
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	if err := migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	log.Printf("Database ready: %s", dbPath)
	return nil
}

func migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS clients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	CREATE TABLE IF NOT EXISTS simulators (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	CREATE TABLE IF NOT EXISTS runs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		client_id INTEGER NOT NULL REFERENCES clients(id),
		simulator_id INTEGER NOT NULL REFERENCES simulators(id),
		date_run TEXT NOT NULL,
		total INTEGER NOT NULL DEFAULT 0,
		passed INTEGER NOT NULL DEFAULT 0,
		failed INTEGER NOT NULL DEFAULT 0,
		version TEXT,
		raw_log_path TEXT,
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS tests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		run_id INTEGER NOT NULL REFERENCES runs(id),
		test_name TEXT NOT NULL,
		passed INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS probes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		client_id INTEGER NOT NULL REFERENCES clients(id),
		version TEXT,
		supported INTEGER NOT NULL DEFAULT 0,
		unsupported INTEGER NOT NULL DEFAULT 0,
		total INTEGER NOT NULL DEFAULT 0,
		modules TEXT,
		date_run TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS probe_methods (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		probe_id INTEGER NOT NULL REFERENCES probes(id),
		method TEXT NOT NULL,
		supported INTEGER NOT NULL,
		sample_value TEXT,
		error TEXT
	);
	CREATE TABLE IF NOT EXISTS comparisons (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		simulator TEXT NOT NULL,
		client_a_name TEXT NOT NULL,
		client_b_name TEXT NOT NULL,
		both_pass INTEGER DEFAULT 0,
		a_only INTEGER DEFAULT 0,
		b_only INTEGER DEFAULT 0,
		both_fail INTEGER DEFAULT 0,
		total_matched INTEGER DEFAULT 0,
		date_compared TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS comparison_tests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		comparison_id INTEGER NOT NULL REFERENCES comparisons(id),
		test_name TEXT NOT NULL,
		a_pass INTEGER,
		b_pass INTEGER,
		status TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS gap_matrices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		client_a_id INTEGER NOT NULL REFERENCES clients(id),
		client_b_id INTEGER NOT NULL REFERENCES clients(id),
		version_a TEXT,
		version_b TEXT,
		modules_a TEXT,
		modules_b TEXT,
		total_methods INTEGER DEFAULT 0,
		both_supported INTEGER DEFAULT 0,
		both_unsupported INTEGER DEFAULT 0,
		in_a_not_b INTEGER DEFAULT 0,
		in_b_not_a INTEGER DEFAULT 0,
		date_created TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS gap_matrix_methods (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		matrix_id INTEGER NOT NULL REFERENCES gap_matrices(id),
		method TEXT NOT NULL,
		a_supported INTEGER NOT NULL,
		b_supported INTEGER NOT NULL,
		a_error TEXT,
		b_error TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_tests_run_id ON tests(run_id);
	CREATE INDEX IF NOT EXISTS idx_probe_methods_probe_id ON probe_methods(probe_id);
	CREATE INDEX IF NOT EXISTS idx_comparison_tests_comparison_id ON comparison_tests(comparison_id);
	CREATE INDEX IF NOT EXISTS idx_gap_matrix_methods_matrix_id ON gap_matrix_methods(matrix_id);
	`
	_, err := db.Exec(schema)
	return err
}

func ensureClient(name string) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT id FROM clients WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := db.Exec("INSERT INTO clients (name) VALUES (?)", name)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	return id, err
}

func ensureSimulator(name string) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT id FROM simulators WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := db.Exec("INSERT INTO simulators (name) VALUES (?)", name)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	return id, err
}
