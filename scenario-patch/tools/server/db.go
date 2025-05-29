package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type DB struct {
	sqlite *sql.DB
	sync.Mutex
}

func (db *DB) Exec(sql string, args ...any) (sql.Result, error) {
	db.Lock()
	defer db.Unlock()
	return db.sqlite.Exec(sql, args...)
}
func (db *DB) Query(sql string, args ...any) (*sql.Rows, error) {
	db.Lock()
	defer db.Unlock()
	return db.sqlite.Query(sql, args...)
}
func (db *DB) QueryRow(sql string, args ...any) *sql.Row {
	db.Lock()
	defer db.Unlock()
	return db.sqlite.QueryRow(sql, args...)
}

var (
	TABLE_PLAYER = `
	CREATE TABLE IF NOT EXISTS players (
	  name TEXT PRIMARY KEY,
	  pocket INTEGER NOT NULL DEFAULT 20,
	  bank INTEGER NOT NULL DEFAULT 10
	);`

	TABLE_FRIEND = `
    CREATE TABLE IF NOT EXISTS friends (
        id INTEGER PRIMARY KEY,
        p1 TEXT REFERENCES players(name),
        p2 TEXT REFERENCES players(name),
        UNIQUE(p1, p2)
    );`

	TABLE_SECRET = `
    CREATE TABLE IF NOT EXISTS secrets (
    	user TEXT,
        pass TEXT,
        PRIMARY KEY (user, pass)
    );`
)

type player struct {
	Name   string
	Pocket int
	Bank   int
}

func (p *player) print() string {
	return fmt.Sprintf(`Name : %s, Bank : %d, Pocket : %d`, p.Name, p.Bank, p.Pocket)
}

var jumpLock = map[string]*sync.Mutex{}

// Lock only the player and not the whole db
func (p *player) jump(db *DB) error {
	db.Lock()
	if _, ok := jumpLock[p.Name]; !ok {
		jumpLock[p.Name] = &sync.Mutex{}
	}
	db.Unlock()
	jumpLock[p.Name].Lock()
	defer jumpLock[p.Name].Unlock()
	_, err := db.Exec(`UPDATE players SET pocket = pocket + 1 WHERE name = ?`, p.Name)
	time.Sleep(1 * time.Second)
	return err
}

func (p *player) insert(db *DB) error {
	_, err := db.Exec(`INSERT INTO players (name, pocket, bank) VALUES (?, ?, ?) ON CONFLICT(name) DO NOTHING`, p.Name, p.Pocket, p.Bank)
	return err
}

func (p *player) addFriend(db *DB, friends []string) error {
	stmt := ""
	for _, friend := range friends {
		stmt += fmt.Sprintf("INSERT INTO friends (p1, p2) VALUES ('%s', '%s');", p.Name, friend)
		stmt += fmt.Sprintf("INSERT INTO friends (p1, p2) VALUES ('%s', '%s');", friend, p.Name)
	}
	_, err := db.Exec(stmt)
	return err
}

func (p *player) findFriends(db *DB) ([]player, error) {
	rows, err := db.Query(`
		SELECT p2.name, p2.pocket, p2.bank
		FROM friends f
		JOIN players p1 ON f.p1 = p1.name
		JOIN players p2 ON f.p2 = p2.name
		WHERE p1.name = ?
	`, p.Name)
	if err != nil {
		return []player{}, err
	}
	friends := []player{}
	for rows.Next() {
		var friend player
		if err = rows.Scan(&friend.Name, &friend.Pocket, &friend.Bank); err != nil {
			return []player{}, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

func (p *player) toBank(db *DB) error {
	if _, err := db.Exec("UPDATE players SET bank = bank + pocket WHERE name = ?", p.Name); err != nil {
		return err
	}
	p.Bank += p.Pocket
	if _, err := db.Exec("UPDATE players SET pocket = 0 WHERE name = ?", p.Name); err != nil {
		return err
	}
	p.Pocket = 0
	return nil
}

func (p *player) toPocket(db *DB) error {
	if _, err := db.Exec("UPDATE players SET pocket = pocket + bank WHERE name = ?", p.Name); err != nil {
		return err
	}
	p.Pocket += p.Bank
	if _, err := db.Exec("UPDATE players SET bank = 0 WHERE name = ?", p.Name); err != nil {
		return err
	}
	p.Bank = 0
	return nil
}

func findPlayer(db *DB, name string) (player, error) {
	row := db.QueryRow("SELECT name, pocket, bank FROM players WHERE name = ?", name)
	if row.Err() != nil {
		return player{}, row.Err()
	}
	p := player{}
	row.Scan(&p.Name, &p.Pocket, &p.Bank)
	return p, nil
}

func isAdmin(db *DB, user string, pass string) bool {
	row := db.QueryRow("SELECT 1 FROM secrets WHERE user = ? AND pass = ?", user, pass)
	if row.Err() != nil {
		return false
	}
	exist := 0
	row.Scan(&exist)
	return exist == 1
}

func initDB() *DB {
	// Open an in-memory SQLite database
	d, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		panic(err)
	}
	db := DB{sqlite: d}

	if _, err = db.Exec(TABLE_PLAYER); err != nil {
		panic(err)
	}
	if _, err = db.Exec(TABLE_FRIEND); err != nil {
		panic(err)
	}
	if _, err = db.Exec(TABLE_SECRET); err != nil {
		panic(err)
	}

	players := []player{
		{Name: "noah", Pocket: 10, Bank: 15},
		{Name: "elann", Pocket: 25, Bank: 15},
		{Name: "jimmy", Pocket: 10, Bank: 15},
	}
	for _, player := range players {
		if err = player.insert(&db); err != nil {
			panic(err)
		}
	}

	if _, err = db.Exec("INSERT INTO secrets (user, pass) VALUES ('admin', 'YWRtaW4=') ON CONFLICT DO NOTHING"); err != nil {
		panic(err)
	}

	return &db
}
