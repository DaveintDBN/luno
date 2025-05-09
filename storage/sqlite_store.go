package storage

import (
    "database/sql"
    _ "github.com/glebarez/sqlite"
    "fmt"
    "time"
)

// SQLiteStore persists trades and slices to a local SQLite database.
type SQLiteStore struct {
    db *sql.DB
}

// NewSQLiteStore opens or creates the database at the given path and runs migrations.
func NewSQLiteStore(path string) (*SQLiteStore, error) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, fmt.Errorf("open sqlite db: %w", err)
    }
    if err := runMigrations(db); err != nil {
        db.Close()
        return nil, fmt.Errorf("run migrations: %w", err)
    }
    return &SQLiteStore{db: db}, nil
}

// Close closes the database connection.
func (s *SQLiteStore) Close() error {
    return s.db.Close()
}

// SaveTrade inserts a trade record and returns its generated ID.
func (s *SQLiteStore) SaveTrade(timestamp time.Time, pair, side string, price, volume float64) (int64, error) {
    rs, err := s.db.Exec(`INSERT INTO trades(timestamp, pair, side, price, volume) VALUES (?, ?, ?, ?, ?)`,
        timestamp.Format(time.RFC3339Nano), pair, side, price, volume)
    if err != nil {
        return 0, err
    }
    return rs.LastInsertId()
}

// SaveSlice inserts a slice record linked to a trade.
func (s *SQLiteStore) SaveSlice(tradeID int64, index int, size, weight float64) error {
    _, err := s.db.Exec(`INSERT INTO slices(trade_id, slice_index, size, weight) VALUES (?, ?, ?, ?)`,
        tradeID, index, size, weight)
    return err
}

// Trade represents a persisted trade record.
type Trade struct {
    ID        int64
    Timestamp time.Time
    Pair      string
    Side      string
    Price     float64
    Volume    float64
}

// ListTrades returns all persisted trades ordered by timestamp.
func (s *SQLiteStore) ListTrades() ([]Trade, error) {
    rows, err := s.db.Query(`SELECT id, timestamp, pair, side, price, volume FROM trades ORDER BY timestamp`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var trades []Trade
    for rows.Next() {
        var t Trade
        var ts string
        if err := rows.Scan(&t.ID, &ts, &t.Pair, &t.Side, &t.Price, &t.Volume); err != nil {
            return nil, err
        }
        t.Timestamp, _ = time.Parse(time.RFC3339Nano, ts)
        trades = append(trades, t)
    }
    return trades, nil
}

// SliceRecord represents a persisted slice record.
type SliceRecord struct {
    ID      int64
    TradeID int64
    Index   int
    Size    float64
    Weight  float64
}

// ListSlices returns all slices for a given trade ID ordered by slice index.
func (s *SQLiteStore) ListSlices(tradeID int64) ([]SliceRecord, error) {
    rows, err := s.db.Query(`SELECT id, trade_id, slice_index, size, weight FROM slices WHERE trade_id = ? ORDER BY slice_index`, tradeID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var slices []SliceRecord
    for rows.Next() {
        var sr SliceRecord
        if err := rows.Scan(&sr.ID, &sr.TradeID, &sr.Index, &sr.Size, &sr.Weight); err != nil {
            return nil, err
        }
        slices = append(slices, sr)
    }
    return slices, nil
}

// runMigrations creates the trades and slices tables if they do not exist.
func runMigrations(db *sql.DB) error {
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS trades (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp TEXT,
        pair TEXT,
        side TEXT,
        price REAL,
        volume REAL
    );`)
    if err != nil {
        return err
    }
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS slices (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        trade_id INTEGER,
        slice_index INTEGER,
        size REAL,
        weight REAL,
        FOREIGN KEY(trade_id) REFERENCES trades(id)
    );`)
    return err
}
