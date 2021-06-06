package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/BitCrackers/BitBot/modlog"

	"github.com/BitCrackers/BitBot/config"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	underlying   *sql.DB
	session      *discordgo.Session
	config       *config.Config
	modlog       *modlog.ModLogHandler
	closeJanitor chan struct{}
	closed       bool
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
	Type      int           `json:"type"`
	Reason    string        `json:"reason"`
	Moderator string        `json:"moderator"`
	Length    time.Duration `json:"length"`
	Date      time.Time     `json:"Date"`
}

type UserRecord struct {
	ID       string
	Warnings []Warning
	Ban      Punishment
	Mute     Punishment
}

type ReactionRole struct {
	ID string `json:"-"`
	Emote string `json:"emote"`
	Role string `json:"role"`
	Channel string `json:"channel"`
}

func New(session *discordgo.Session, config *config.Config, modlog *modlog.ModLogHandler) (*Database, error) {
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

	underlying, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error while opening db file: %v", err)
	}

	if sqlQuery != "" {
		_, err = underlying.Exec(sqlQuery)
		if err != nil {
			return nil, err
		}
		sqlQuery = `
	CREATE TABLE reactionroles (
	id VARCHAR NOT NULL,
	data LONGTEXT
);
	`
		_, err = underlying.Exec(sqlQuery)
		if err != nil {
			return nil, err
		}
	}
	db := &Database{
		underlying: underlying,
		session:    session,
		config:     config,
		modlog:     modlog,
	}
	db.closeJanitor = db.startJanitor()
	return db, nil
}

func (d *Database) WarnUser(id string, moderatorID string, reason string) error {
	userRecord, err := d.UserRecord(id)
	if err != nil {
		return err
	}
	userRecord.Warnings = append(userRecord.Warnings, Warning{
		Reason:    reason,
		Moderator: moderatorID,
		Date:      time.Now(),
	})
	return d.SetUserRecord(userRecord)
}

func (d *Database) MuteUser(id string, moderatorID string, reason string, length time.Duration) error {
	record, err := d.UserRecord(id)
	if err != nil {
		return fmt.Errorf("error while getting user record: %v", err)
	}
	if !record.Mute.Empty() {
		remaining := ""
		if record.Mute.Length != -1 {
			remaining = fmt.Sprintf(
				", %v remaining",
				record.Mute.Date.Add(record.Mute.Length).Sub(time.Now()).String(),
			)
		}
		return fmt.Errorf("user has already been muted%v", remaining)
	}

	err = d.session.GuildMemberRoleAdd(d.config.GuildID, id, d.config.MuteRoleID)
	if err != nil {
		return fmt.Errorf("error while adding muted role to user: %v", err)
	}

	return d.PunishUser(id, Punishment{
		Type:      PunishmentTypeMute,
		Reason:    reason,
		Moderator: moderatorID,
		Length:    length,
		Date:      time.Now(),
	})
}

func (d *Database) BanUser(id string, moderatorID string, reason string, length time.Duration) error {
	if err := d.session.GuildBanCreateWithReason(d.config.GuildID, id, reason, 0); err != nil {
		return fmt.Errorf("error while banning user from guild: %v", err)
	}
	return d.PunishUser(id, Punishment{
		Type:      PunishmentTypeBan,
		Reason:    reason,
		Moderator: moderatorID,
		Length:    length,
		Date:      time.Now(),
	})
}

func (d *Database) PunishUser(id string, punishment Punishment) error {
	userRecord, err := d.UserRecord(id)
	if err != nil {
		return err
	}

	if punishment.Type == PunishmentTypeBan {
		userRecord.Ban = punishment
	}
	if punishment.Type == PunishmentTypeMute {
		userRecord.Mute = punishment
	}

	if err = d.SetUserRecord(userRecord); err != nil {
		return err
	}
	return nil
}

func (d *Database) UserRecordExists(id string) (bool, error) {
	sqlStmt := `SELECT id FROM userinfo WHERE id = ?`
	err := d.underlying.QueryRow(sqlStmt, id).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (d *Database) UserRecord(id string) (UserRecord, error) {
	e, err := d.UserRecordExists(id)
	if err != nil {
		return UserRecord{}, err
	}

	if !e {
		err = d.CreateUserRecord(id)
		if err != nil {
			return UserRecord{}, err
		}
	}

	sqlStmt := `SELECT id, warnings, mute, ban FROM userinfo WHERE id = ?`
	row := d.underlying.QueryRow(sqlStmt, id)

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

	if err = row.Scan(&idNString, &warningsNString, &muteNString, &banNString); err != nil {
		return UserRecord{}, err
	}
	userRecord.ID = idNString.String

	warningsString := ""
	if warningsNString.Valid {
		warningsString = warningsNString.String
	}
	if warningsString != "" {
		if err = json.Unmarshal([]byte(warningsString), &userRecord.Warnings); err != nil {
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
		if err = json.Unmarshal([]byte(muteString), &mute); err != nil {
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
		if err = json.Unmarshal([]byte(banString), &ban); err != nil {
			return UserRecord{}, err
		}
	}
	userRecord.Ban = ban

	return userRecord, nil
}

func (d *Database) AllUserRecords() ([]UserRecord, error) {
	sqlStmt := `SELECT id, warnings, mute, ban FROM userinfo`
	rows, err := d.underlying.Query(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error while fetching user records for database: %v", err)
	}

	var records []UserRecord
	for rows.Next() {
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

		if err = rows.Scan(&idNString, &warningsNString, &muteNString, &banNString); err != nil {
			return nil, err
		}
		userRecord.ID = idNString.String

		warningsString := ""
		if warningsNString.Valid {
			warningsString = warningsNString.String
		}
		if warningsString != "" {
			if err = json.Unmarshal([]byte(warningsString), &userRecord.Warnings); err != nil {
				return nil, err
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
			if err = json.Unmarshal([]byte(muteString), &mute); err != nil {
				return nil, err
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
			if err = json.Unmarshal([]byte(banString), &ban); err != nil {
				return nil, err
			}
		}
		userRecord.Ban = ban

		records = append(records, userRecord)
	}
	return records, nil
}

func (d *Database) SetUserRecord(record UserRecord) error {
	e, err := d.UserRecordExists(record.ID)
	if err != nil {
		return err
	}

	if !e {
		if err = d.CreateUserRecord(record.ID); err != nil {
			return err
		}
	}

	warnings := sql.NullString{
		String: "",
		Valid:  false,
	}
	if len(record.Warnings) > 0 {
		b, err := json.Marshal(record.Warnings)
		if err != nil {
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

func (d *Database) CreateUserRecord(id string) error {
	tx, err := d.underlying.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO userinfo (id) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) AddReactionRole(id string, reaction ReactionRole) error  {
	tx, err := d.underlying.Begin()
	if err != nil {
		return err
	}

	r := sql.NullString{
		String: "",
		Valid:  true,
	}

	b, err := json.Marshal(reaction)
	if err != nil {
		return err
	}
	r.String = string(b)

	stmt, err := tx.Prepare("INSERT INTO reactionroles (id, data) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, r)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) ReactionRoleExist(id string) (bool, error)  {
	sqlStmt := `SELECT id FROM reactionroles WHERE id = ?`
	err := d.underlying.QueryRow(sqlStmt, id).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (d *Database) GetReaction(id string) (ReactionRole, error) {
	e, err := d.ReactionRoleExist(id)
	if err != nil {
		return ReactionRole{}, err
	}

	if !e {
		return ReactionRole{}, nil
	}

	sqlStmt := `SELECT id, data FROM reactionroles WHERE id = ?`
	row := d.underlying.QueryRow(sqlStmt, id)

	var idNString sql.NullString
	var dataNString sql.NullString

	if err = row.Scan(&idNString, &dataNString); err != nil {
		return ReactionRole{}, err
	}

	r := ReactionRole{}

	dataString := ""
	if dataNString.Valid {
		dataString = dataNString.String
	}
	if dataString != "" {
		if err = json.Unmarshal([]byte(dataString), &r); err != nil {
			return ReactionRole{}, err
		}
	}
	r.ID = idNString.String

	return r, nil
}

func (d *Database) AllReactionRole() ([]ReactionRole, error) {
	sqlStmt := `SELECT id, data FROM reactionroles`
	rows, err := d.underlying.Query(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error while fetching user records for database: %v", err)
	}

	var records []ReactionRole
	for rows.Next() {
		var idNString sql.NullString
		var dataNString sql.NullString

		if err = rows.Scan(&idNString, &dataNString); err != nil {
			return nil, err
		}

		r := ReactionRole{}

		dataString := ""
		if dataNString.Valid {
			dataString = dataNString.String
		}
		if dataString != "" {
			if err = json.Unmarshal([]byte(dataString), &r); err != nil {
				return nil, err
			}
		}
		r.ID = idNString.String

		records = append(records, r)
	}
	return records, nil
}

func (d *Database) RemoveRoleReaction(id string) error  {
	tx, err := d.underlying.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("DELETE FROM reactionroles WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
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
	if d.closed {
		panic("attempted to close an already closed database")
	}
	d.closeJanitor <- struct{}{}
	close(d.closeJanitor)
	d.closed = true
	return d.underlying.Close()
}

func (p *Punishment) Empty() bool {
	return p.Type == -1
}
