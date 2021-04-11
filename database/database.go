package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Database struct {
	underlying *sql.DB
}

const (
	PunishmentTypeMute = iota
	PunishmentTypeBan
)

type Warning struct {
	Reason    string    `json:"reason"`
	Moderator string    `json:"moderator"`
	Date      time.Time `json:"date"`
}

type Punishment struct {
	Type      int       `json:"type"`
	Reason    string    `json:"reason"`
	Moderator string    `json:"moderator"`
	Length    int       `json:"length"`
	Date      time.Time `json:"Date"`
}

func New() (*Database, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("error while getting executable directory: %v", err)
	}
	ex = filepath.ToSlash(ex) // For the operating systems which use backslash...
	dbPath := path.Join(path.Dir(ex), "bitbot.db")

	sqlQuery := ""
	if _, err = os.Stat(dbPath); os.IsNotExist(err) {
		sqlQuery = `
	CREATE TABLE userinfo (
	id INT NOT NULL,
	warnings LONGTEXT,
	mute LONGTEXT,
	ban LONGTEXT
);
	`
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error while opening db file: %v", err)
	}

	if sqlQuery != "" {
		_, err = db.Exec(sqlQuery)
		if err != nil {
			return nil, err
		}
	}

	return &Database{db}, nil
}

func (d *Database) WarnUser(user *discordgo.User, moderator *discordgo.User, reason string) error {
	e, err := d.UserRecordExists(user)
	if err != nil {
		return err
	}

	if !e {
		err = d.CreateUserRecord(user)
		if err != nil {
			return err
		}
	}

	sqlStmt := `SELECT warnings FROM userinfo WHERE id = ?`
	row := d.underlying.QueryRow(sqlStmt, user.ID)
	var s sql.NullString

	err = row.Scan(&s)
	if err != nil {
		return err
	}

	warningsString := ""
	if s.Valid {
		warningsString = s.String
	}

	var warnings []Warning
	if warningsString != "" {
		err = json.Unmarshal([]byte(warningsString), &warnings)
		if err != nil {
			fmt.Printf("unable to unmarshal warnings object: %s\n", warningsString)
			return err
		}
	}

	if reason == "" {
		reason = "unknown"
	}

	warnings = append(warnings, Warning{
		Reason:    reason,
		Moderator: moderator.ID,
		Date:      time.Now(),
	})

	b, err := json.Marshal(warnings)
	if err != nil {
		fmt.Printf("unable to marshal warnings object\n")
		return err
	}

	tx, err := d.underlying.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE userinfo SET warnings = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(string(b), user.ID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) UserRecordExists(user *discordgo.User) (bool, error) {
	sqlStmt := `SELECT id FROM userinfo WHERE id = ?`
	var id int
	err := d.underlying.QueryRow(sqlStmt, user.ID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (d *Database) CreateUserRecord(user *discordgo.User) error {
	tx, err := d.underlying.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO userinfo (id) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) Close() error {
	return d.underlying.Close()
}
