package model

import (
	"github.com/twinj/uuid"
	"log"
	"strings"
	"time"
)

type Tag struct {
	Id        int64
	Name      string
	Slug      string
	CreatedAt *time.Time
	CreatedBy int64
	Hidden    bool
}

func (t *Tag) Url() string {
	return "/tag/" + t.Slug
}

func (t *Tag) Save() error {
	oldTag, err := GetTagBySlug(t.Slug)
	if err != nil {
		// Tag is probably not in database yet

		tagId, err := InsertTag(t)
		if err != nil {
			log.Printf("[Error] Can not insert tag: %v", err.Error())
			return err
		}
		t.Id = tagId
	} else {
		t.Id = oldTag.Id
		// If oldTag.Hidden != t.Hidden, we need to decide whether change the hidden status of oldTag
		if oldTag.Hidden != t.Hidden {
			if oldTag.Hidden {
				err = UpdateTag(t)
				if err != nil {
					return err
				}
			} else {
				// oldTag.Hidden is false and t.Hidden is true
				posts, err := GetAllPostsByTag(oldTag.Id)
				needUpdate := true
				for _, p := range posts {
					if p.IsPublished {
						needUpdate = false
						break
					}
				}
				if needUpdate {
					err = UpdateTag(t)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func GenerateTagsFromCommaString(input string) []*Tag {
	output := make([]*Tag, 0)
	tags := strings.Split(input, ",")
	for index, _ := range tags {
		tags[index] = strings.TrimSpace(tags[index])
	}
	for _, tag := range tags {
		if tag != "" {
			output = append(output, &Tag{Name: tag, Slug: GenerateSlug(tag, "tags")})
		}
	}
	return output
}

func GetTagIdBySlug(slug string) (int64, error) {
	var id int64
	row := db.QueryRow(stmtGetTagIdBySlug, slug)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func InsertTag(t *Tag) (int64, error) {
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	result, err := writeDB.Exec(stmtInsertTag, nil, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), t.Name, t.Slug, t.CreatedAt, t.CreatedBy, t.CreatedAt, t.CreatedBy, t.Hidden)
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	tagId, err := result.LastInsertId()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	return tagId, writeDB.Commit()
}

func UpdateTag(t *Tag) error {
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateTag, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), t.Name, t.Slug, t.CreatedAt, t.CreatedBy, t.Hidden, t.Id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func GetTags(postId int64) ([]*Tag, error) {
	tags := make([]*Tag, 0)
	// Get tags
	rows, err := db.Query(stmtGetTags, postId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tagId int64
		err := rows.Scan(&tagId)
		if err != nil {
			return nil, err
		}
		tag, err := GetTag(tagId)
		// TODO: Error while receiving individual tag is ignored right now. Keep it this way?
		if err == nil {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

func GetTag(tagId int64) (*Tag, error) {
	tag := new(Tag)
	// Get tag
	row := db.QueryRow(stmtGetTagById, tagId)
	err := row.Scan(&tag.Id, &tag.Name, &tag.Slug)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func GetTagBySlug(slug string) (*Tag, error) {
	tag := new(Tag)
	// Get tag
	row := db.QueryRow(stmtGetTagBySlug, slug)
	err := row.Scan(&tag.Id, &tag.Name, &tag.Slug, &tag.Hidden)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func GetAllTags() ([]*Tag, error) {
	tags := make([]*Tag, 0)
	// Get tags
	rows, err := db.Query(stmtGetAllTags)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		tag := new(Tag)
		err := rows.Scan(&tag.Id, &tag.Name, &tag.Slug)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func CleanTags() error {
	return nil
}

func DeleteOldTags() error {
	WriteDB, err := db.Begin()
	if err != nil {
		WriteDB.Rollback()
		return err
	}
	_, err = WriteDB.Exec(stmtDeleteOldTags)
	if err != nil {
		WriteDB.Rollback()
		return err
	}
	return WriteDB.Commit()
}
