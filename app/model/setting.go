package model

import (
	"database/sql"
	"encoding/json"
	"github.com/twinj/uuid"
	"time"
)

type Setting struct {
	Id        int
	UUID      string
	Key       string
	Value     string
	Type      string // general, content, navigation, custom
	CreatedAt *time.Time
	CreatedBy int64
}

type Navigator struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

func GetNavigators() []*Navigator {
	var navs []*Navigator
	navStr := GetSettingValue("navigation")
	json.Unmarshal([]byte(navStr), &navs)
	return navs
}

func SetNavigators(labels, urls []string) error {
	var navs []*Navigator
	for i, l := range labels {
		if len(l) < 1 {
			continue
		}
		navs = append(navs, &Navigator{l, urls[i]})
	}
	navStr, err := json.Marshal(navs)
	if err != nil {
		return err
	}
	err = SetSetting("navigation", string(navStr), "navigation")
	return err
}

func GetSetting(k string) (*Setting, error) {
	row := db.QueryRow(`SELECT id, uuid, key, value, type, created_at, created_by from settings where key = ?`, k)
	s := new(Setting)
	err := scanSetting(row, s)
	return s, err
}

func GetSettingValue(k string) string {
	row := db.QueryRow(`SELECT id, uuid, key, value, type, created_at, created_by from settings where key = ?`, k)
	s := new(Setting)
	scanSetting(row, s)
	return s.Value
}

func scanSetting(row Row, setting *Setting) error {
	var nullValue sql.NullString
	err := row.Scan(&setting.Id, &setting.UUID, &setting.Key, &nullValue, &setting.Type, &setting.CreatedAt, &setting.CreatedBy)
	if err != nil {
		return err
	}
	setting.Value = nullValue.String
	return nil
}

func GetCustomSettings() []*Setting {
	return GetSettings("custom")
}

func GetSettings(t string) []*Setting {
	settings := make([]*Setting, 0)
	rows, err := db.Query(`SELECT id, uuid, key, value, type, created_at, created_by from settings where type = ?`, t)
	if err != nil {
		return settings
	}
	defer rows.Close()
	if err != nil {
		return settings
	}
	for rows.Next() {
		setting := new(Setting)
		err := scanSetting(rows, setting)
		if err != nil {
			return settings
		}
		settings = append(settings, setting)
	}
	return settings

}

func SaveSetting(setting *Setting) error {
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = db.Exec(`INSERT OR REPLACE INTO settings (id, uuid, key, value, type, created_at, created_by) VALUES ((SELECT id FROM settings WHERE key = ?), ?, ?, ?, ?, ?, ?)`, setting.Key, setting.UUID, setting.Key, setting.Value, setting.Type, setting.CreatedAt, setting.CreatedBy)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func SetSetting(k, v, t string) error {
	setting := new(Setting)
	setting.UUID = uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)
	setting.Key = k
	setting.Value = v
	setting.Type = t
	now := time.Now()
	setting.CreatedAt = &now
	return SaveSetting(setting)
}

func SetSettingIfNotExists(k, v, t string) error {
	_, err := GetSetting(k)
	if err != nil {
		err = SetSetting(k, v, t)
	}
	return err
}
