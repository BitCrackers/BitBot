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

type UserRecord struct {
	ID       string
	Warnings []Warning
	Ban      Punishment
	Mute     Punishment
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
	id VARCHAR NOT NULL,
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
	userRecord, err := d.GetUserRecord(user)
	if err != nil {
		return err
	}
	userRecord.Warnings = append(userRecord.Warnings, Warning{
		Reason:    reason,
		Moderator: moderator.ID,
		Date:      time.Now(),
	})
	err = d.SetUserRecord(userRecord)
	return err
}

func (d *Database) MuteUser(user *discordgo.User, moderator *discordgo.User, reason string, length int) error {
	return d.PunishUser(user, Punishment{
		Type:      PunishmentTypeMute,
		Reason:    reason,
		Moderator: moderator.ID,
		Length:    length,
		Date:      time.Now(),
	})
}

func (d *Database) BanUser(user *discordgo.User, moderator *discordgo.User, reason string, length int) error {
	return d.PunishUser(user, Punishment{
		Type:      PunishmentTypeBan,
		Reason:    reason,
		Moderator: moderator.ID,
		Length:    length,
		Date:      time.Now(),
	})
}

func (d *Database) PunishUser(user *discordgo.User, punishment Punishment) error {
	userRecord, err := d.GetUserRecord(user)
	if err != nil {
		return err
	}

	if punishment.Type == PunishmentTypeBan {
		userRecord.Ban = punishment
	}
	if punishment.Type == PunishmentTypeMute {
		userRecord.Mute = punishment
	}

	err = d.SetUserRecord(userRecord)
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

func (d *Database) GetUserRecord(user *discordgo.User) (UserRecord, error) {
	e, err := d.UserRecordExists(user)
	if err != nil {
		return UserRecord{}, err
	}

	if !e {
		err = d.CreateUserRecord(user)
		if err != nil {
			return UserRecord{}, err
		}
	}

	sqlStmt := `SELECT id, warnings, mute, ban FROM userinfo WHERE id = ?`
	row := d.underlying.QueryRow(sqlStmt, user.ID)

	var idNString sql.NullString
	var warningsNString sql.NullString
	var muteNString sql.NullString
	var banNString sql.NullString

	userRecord := UserRecord{
		ID:       "",
		Warnings: []Warning{},
		Ban:      Punishment{},
		Mute:     Punishment{},
	}

	err = row.Scan(&idNString, &warningsNString, &muteNString, &banNString)
	if err != nil {
		return UserRecord{}, err
	}
	userRecord.ID = idNString.String

	warningsString := ""
	if warningsNString.Valid {
		warningsString = warningsNString.String
	}
	if warningsString != "" {
		err = json.Unmarshal([]byte(warningsString), &userRecord.Warnings)
		if err != nil {
			fmt.Printf("unable to unmarshal warnings object: %s\n", warningsString)
			return UserRecord{}, err
		}
	}

	muteString := ""
	if muteNString.Valid {
		muteString = muteNString.String
	}
	mute := Punishment{
		Type: -1,
	}
	if muteString != "" {
		err = json.Unmarshal([]byte(muteString), &mute)
		if err != nil {
			fmt.Printf("unable to unmarshal mute object: %s", muteString)
			return UserRecord{}, err
		}
	}
	userRecord.Mute = mute

	banString := ""
	if banNString.Valid {
		banString = banNString.String
	}
	ban := Punishment{
		Type: -1,
	}
	if banString != "" {
		err = json.Unmarshal([]byte(banString), &ban)
		if err != nil {
			fmt.Printf("unable to unmarshal ban object: %s\n", banString)
			return UserRecord{}, err
		}
	}
	userRecord.Ban = ban

	return userRecord, nil
}

func (d *Database) SetUserRecord(record UserRecord) error {
	warnings := sql.NullString{
		String: "",
		Valid:  false,
	}
	if len(record.Warnings) > 0 {
		b, err := json.Marshal(record.Warnings)
		if err != nil {
			fmt.Printf("unable to marshal warnings object\n")
			return err
		}
		warnings.String = string(b)
		warnings.Valid = true
	}

	mute := sql.NullString{
		String: "",
		Valid:  false,
	}
	if !record.Mute.Empty() {
		b, err := json.Marshal(record.Mute)
		if err != nil {
			fmt.Printf("unable to marshal mute object\n")
			return err
		}
		mute.String = string(b)
		mute.Valid = true
	}

	ban := sql.NullString{
		String: "",
		Valid:  false,
	}
	if !record.Ban.Empty() {
		b, err := json.Marshal(record.Mute)
		if err != nil {
			fmt.Printf("unable to marshal ban object\n")
			return err
		}
		ban.String = string(b)
		ban.Valid = true
	}

	tx, err := d.underlying.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE userinfo SET warnings = ?, ban = ?, mute = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(warnings, ban, mute, record.ID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
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

func (p *Punishment) Empty() bool {
	return p.Type == -1
}