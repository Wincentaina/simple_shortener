package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/omeid/pgerror"
	"good_shortener/internal/storage"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(dbname, passwd string, port int) (*Storage, error) {
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf(
		"user=postgres password=%v dbname=%v sslmode=disable port=%v", passwd, dbname, port,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS url(
		id SERIAL PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		blocked BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP
	    );
	CREATE INDEX IF NOT EXISTS idx ON url(alias);
	CREATE TABLE IF NOT EXISTS "user"(
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		passwordHash TEXT NOT NULL,
		isAdmin BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP
	);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string, userId int64) error {
	const op = "storage.postgres.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias, user_id, created_at) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(urlToSave, alias, userId, time.Now())
	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	stmt, err := s.db.Prepare(`SELECT url FROM url where alias = $1`)
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var res string
	err = stmt.QueryRow(alias).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return res, nil
}

// TODO: DELETE URL
// func (s *Storage) DeleteURL(alias string) error

type User struct {
	Id           int64
	Username     string
	PasswordHash string
	IsAdmin      bool
}

// TODO: remove hardcoded value isAdmin
// returns id, error
func (s *Storage) SaveUser(username, passwordHash string, isAdmin bool) error {
	const op = "storage.postgres.SaveUser"
	isAdmin = false

	stmt, err := s.db.Prepare(`
		INSERT INTO "user"(username, passwordhash, isadmin, created_at) VALUES($1, $2, $3, $4);
	`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	_, err = stmt.Exec(username, passwordHash, isAdmin, time.Now())
	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return fmt.Errorf("%s: %w", op, storage.UserExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetUserById(id int64) (User, error) {
	const op = "storage.postgres.GetUserById"

	stmt, err := s.db.Prepare(`
		SELECT id, username, passwordhash, isadmin FROM "user" WHERE id = $1 LIMIT 1;
	`)
	if err != nil {
		return User{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var usr User
	err = stmt.QueryRow(id).Scan(&usr.Id, &usr.Username, &usr.PasswordHash, &usr.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, storage.UserNotFound
		}

		return User{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return usr, nil
}

func (s *Storage) GetUserByUsername(username string) (User, error) {
	const op = "storage.postgres.GetUserByUsername"

	stmt, err := s.db.Prepare(`
		SELECT id, username, passwordhash, isadmin FROM "user" WHERE username = $1 LIMIT 1;
	`)
	if err != nil {
		return User{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var usr User
	err = stmt.QueryRow(username).Scan(&usr.Id, &usr.Username, &usr.PasswordHash, &usr.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, storage.UserNotFound
		}

		return User{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return usr, nil
}

type Url struct {
	Alias string
	Url   string
}

// TODO: дописать логику
func (s *Storage) GetUserUrls(user_id int64) ([]Url, error) {
	const op = "storage.postgres.GetUserUrls"

	rows, err := s.db.Query(`
		SELECT alias, url FROM "url" WHERE user_id = $1;
	`, user_id)
	if err != nil {
		return []Url{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer rows.Close()

	var urls []Url
	var url Url
	for rows.Next() {
		err := rows.Scan(&url.Alias, &url.Url)
		urls = append(urls, url)
		if err != nil {
			return []Url{}, fmt.Errorf("%s: execute statement: %w", op, err)
		}
	}
	err = rows.Err()
	if err != nil {
		return []Url{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return urls, nil
}
