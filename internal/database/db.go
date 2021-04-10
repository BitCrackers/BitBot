package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path"
	"time"
)

type Database struct {
	db *sql.DB
}

const (
	Mute = iota
	Ban
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

var DB Database

func Start() error {
	wd, err := os.Getwd()
	DB = Database{}

	if err != nil {
		return err
	}
	sqlQuery := ``
	if _, err = os.Stat(path.Join(wd, "bitbot.db")); os.IsNotExist(err) {
		sqlQuery = `
	CREATE TABLE userinfo (
	id INT NOT NULL,
	warnings LONGTEXT,
	mute LONGTEXT,
	ban LONGTEXT
);
	`
	}

	db, err := sql.Open("sqlite3", path.Join(wd, "bitbot.db"))
	if err != nil {
		return err
	}

	if sqlQuery != "" {
		_, err = db.Exec(sqlQuery)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlQuery)
			return err
		}
	}

	DB.db = db
	return nil
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
	row := d.db.QueryRow(sqlStmt, user.ID)

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

	tx, err := d.db.Begin()
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
	err := d.db.QueryRow(sqlStmt, user.ID).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func (d *Database) CreateUserRecord(user *discordgo.User) error {
	tx, err := d.db.Begin()
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
	return d.db.Close()
}
