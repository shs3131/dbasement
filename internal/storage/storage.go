package storage

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	db        *sql.DB
	mu        sync.RWMutex
	projectID int
}

type Section struct {
	ID        int
	Section   string
	Subkey    string
	Content   string
	UpdatedAt time.Time
}

type ChangelogEntry struct {
	ID        int
	Timestamp time.Time
	Entry     string
}

type DesignDecision struct {
	ID        int
	Timestamp time.Time
	Decision  string
	Reason    string
}

type TodoItem struct {
	ID        int
	Timestamp time.Time
	Item      string
	Source    string
	Done      bool
}

type KnownIssue struct {
	ID         int
	Timestamp  time.Time
	Issue      string
	Confidence int
	Resolved   bool
}

func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	d := &DB{db: db}
	if err := d.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return d, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) migrate() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS meta (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sections (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			section TEXT NOT NULL,
			subkey TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			UNIQUE(section, subkey)
		)`,
		`CREATE TABLE IF NOT EXISTS changelog (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL DEFAULT (datetime('now')),
			entry TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS design_decisions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL DEFAULT (datetime('now')),
			decision TEXT NOT NULL,
			reason TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS todo (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL DEFAULT (datetime('now')),
			item TEXT NOT NULL,
			source TEXT NOT NULL DEFAULT 'ai',
			done INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS known_issues (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL DEFAULT (datetime('now')),
			issue TEXT NOT NULL,
			confidence INTEGER NOT NULL DEFAULT 100,
			resolved INTEGER NOT NULL DEFAULT 0
		)`,
	}

	for _, stmt := range statements {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("exec %q: %w", stmt[:40], err)
		}
	}

	return tx.Commit()
}

func (d *DB) SetMeta(key, value string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(
		`INSERT INTO meta (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	return err
}

func (d *DB) GetMeta(key string) (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var value string
	err := d.db.QueryRow(`SELECT value FROM meta WHERE key = ?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (d *DB) SetSection(section, subkey, content string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(
		`INSERT INTO sections (section, subkey, content, updated_at)
		VALUES (?, ?, ?, datetime('now'))
		ON CONFLICT(section, subkey) DO UPDATE SET
			content = excluded.content,
			updated_at = excluded.updated_at`,
		section, subkey, content,
	)
	return err
}

func (d *DB) GetSection(section, subkey string) (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var content string
	err := d.db.QueryRow(
		`SELECT content FROM sections WHERE section = ? AND subkey = ?`,
		section, subkey,
	).Scan(&content)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return content, err
}

func (d *DB) GetSectionAll(section string) ([]Section, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(
		`SELECT id, section, subkey, content, updated_at FROM sections WHERE section = ? ORDER BY id`,
		section,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Section
	for rows.Next() {
		var s Section
		var updated string
		if err := rows.Scan(&s.ID, &s.Section, &s.Subkey, &s.Content, &updated); err != nil {
			return nil, err
		}
		s.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updated)
		result = append(result, s)
	}
	return result, rows.Err()
}

func (d *DB) DeleteSection(section, subkey string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(
		`DELETE FROM sections WHERE section = ? AND subkey = ?`,
		section, subkey,
	)
	return err
}

func (d *DB) AddChangelog(entry string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(
		`INSERT INTO changelog (timestamp, entry) VALUES (datetime('now'), ?)`,
		entry,
	)
	return err
}

func (d *DB) GetChangelog(limit int) ([]ChangelogEntry, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if limit <= 0 {
		limit = 50
	}

	rows, err := d.db.Query(
		`SELECT id, timestamp, entry FROM changelog ORDER BY id DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ChangelogEntry
	for rows.Next() {
		var e ChangelogEntry
		var ts string
		if err := rows.Scan(&e.ID, &ts, &e.Entry); err != nil {
			return nil, err
		}
		e.Timestamp, _ = time.Parse("2006-01-02 15:04:05", ts)
		result = append(result, e)
	}
	return result, rows.Err()
}

func (d *DB) AddDesignDecision(decision, reason string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(
		`INSERT INTO design_decisions (timestamp, decision, reason) VALUES (datetime('now'), ?, ?)`,
		decision, reason,
	)
	return err
}

func (d *DB) GetDesignDecisions() ([]DesignDecision, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(
		`SELECT id, timestamp, decision, reason FROM design_decisions ORDER BY id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []DesignDecision
	for rows.Next() {
		var dd DesignDecision
		var ts string
		if err := rows.Scan(&dd.ID, &ts, &dd.Decision, &dd.Reason); err != nil {
			return nil, err
		}
		dd.Timestamp, _ = time.Parse("2006-01-02 15:04:05", ts)
		result = append(result, dd)
	}
	return result, rows.Err()
}

func (d *DB) AddTodo(item, source string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(
		`INSERT INTO todo (timestamp, item, source) VALUES (datetime('now'), ?, ?)`,
		item, source,
	)
	return err
}

func (d *DB) GetTodos(includeDone bool) ([]TodoItem, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	query := `SELECT id, timestamp, item, source, done FROM todo`
	if !includeDone {
		query += ` WHERE done = 0`
	}
	query += ` ORDER BY id DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TodoItem
	for rows.Next() {
		var t TodoItem
		var ts string
		if err := rows.Scan(&t.ID, &ts, &t.Item, &t.Source, &t.Done); err != nil {
			return nil, err
		}
		t.Timestamp, _ = time.Parse("2006-01-02 15:04:05", ts)
		result = append(result, t)
	}
	return result, rows.Err()
}

func (d *DB) MarkTodoDone(id int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`UPDATE todo SET done = 1 WHERE id = ?`, id)
	return err
}

func (d *DB) AddKnownIssue(issue string, confidence int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(
		`INSERT INTO known_issues (timestamp, issue, confidence) VALUES (datetime('now'), ?, ?)`,
		issue, confidence,
	)
	return err
}

func (d *DB) GetKnownIssues(includeResolved bool) ([]KnownIssue, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	query := `SELECT id, timestamp, issue, confidence, resolved FROM known_issues`
	if !includeResolved {
		query += ` WHERE resolved = 0`
	}
	query += ` ORDER BY id DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []KnownIssue
	for rows.Next() {
		var ki KnownIssue
		var ts string
		if err := rows.Scan(&ki.ID, &ts, &ki.Issue, &ki.Confidence, &ki.Resolved); err != nil {
			return nil, err
		}
		ki.Timestamp, _ = time.Parse("2006-01-02 15:04:05", ts)
		result = append(result, ki)
	}
	return result, rows.Err()
}

func (d *DB) ResolveKnownIssue(id int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`UPDATE known_issues SET resolved = 1 WHERE id = ?`, id)
	return err
}

func (d *DB) SearchMemory(query string) (map[string][]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	results := make(map[string][]string)
	like := "%" + query + "%"

	sections, err := d.db.Query(
		`SELECT section, subkey, content FROM sections WHERE content LIKE ? OR section LIKE ? OR subkey LIKE ?`,
		like, like, like,
	)
	if err != nil {
		return nil, err
	}
	defer sections.Close()

	for sections.Next() {
		var sec, subkey, content string
		if err := sections.Scan(&sec, &subkey, &content); err != nil {
			return nil, err
		}
		label := sec
		if subkey != "" {
			label = sec + "/" + subkey
		}
		results[label] = append(results[label], content)
	}

	return results, nil
}
