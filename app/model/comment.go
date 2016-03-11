package model

import (
	"database/sql"
	"github.com/dinever/dingo/app/utils"
	"github.com/twinj/uuid"
	"time"
)

// Comment struct defines a comment item data.
type Comment struct {
	Id        int64
	UUID      string
	Author    string
	Email     string
	Website   string
	Avatar    string
	CreatedAt *time.Time
	Content   string
	Approved  bool
	PostId    int64
	Parent    int64
	Type      string
	Ip        string
	UserAgent string
	UserId    int64
}

func (comment *Comment) Save() (int64, error) {
	createdAt := time.Now()
	comment.CreatedAt = &createdAt
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	var result sql.Result
	if comment.Id > 0 {
		result, err = writeDB.Exec(stmtInsertComment, comment.Id, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), comment.PostId, comment.Author, comment.Email, comment.Website, comment.Ip, comment.CreatedAt, comment.Content, comment.Approved, comment.UserAgent, comment.Parent, comment.UserId)
	} else {
		result, err = writeDB.Exec(stmtInsertComment, nil, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), comment.PostId, comment.Author, comment.Email, comment.Website, comment.Ip, comment.CreatedAt, comment.Content, comment.Approved, comment.UserAgent, comment.Parent, comment.UserId)
	}
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	commentId, err := result.LastInsertId()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	comment.Id = commentId
	return commentId, writeDB.Commit()
}

func (c *Comment) ToJson() map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = c.Id
	m["author"] = c.Author
	m["email"] = c.Email
	m["website"] = c.Website
	m["avatar"] = c.Avatar
	m["content"] = c.Content
	m["create_time"] = c.CreatedAt.Unix()
	m["pid"] = c.Parent
	m["approved"] = c.Approved
	m["ip"] = c.Ip
	m["user_agent"] = c.UserAgent
	m["parent_content"] = c.ParentContent()
	return m
}

func (c *Comment) ParentContent() string {
	if c.Parent < 1 {
		return ""
	}
	comment, err := GetCommentById(c.Parent)
	if err != nil {
		return "> Comment not found."
	}
	str := "> @" + comment.Author + "\n\n"
	str += "> " + comment.Content + "\n"
	return str
}

func GetNumberOfComments() (int64, error) {
	var count int64
	var row *sql.Row
	row = db.QueryRow(stmtGetAllCommentCount)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetCommentList(page, size int64) ([]*Comment, *utils.Pager, error) {
	var (
		pager *utils.Pager
	)
	count, err := GetNumberOfComments()
	pager = utils.NewPager(page, size, count)
	rows, err := db.Query(stmtGetAllCommentList, size, pager.Begin-1)
	defer rows.Close()
	if err != nil {
		return nil, nil, err
	}
	comments, err := extractComments(rows)
	if err != nil {
		return nil, nil, err
	}
	return comments, pager, nil
}

func extractComments(rows *sql.Rows) ([]*Comment, error) {
	comments := make([]*Comment, 0)
	for rows.Next() {
		comment := new(Comment)
		err := scanComment(rows, comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func scanComment(rows Row, comment *Comment) error {
	var (
		nullParent sql.NullInt64
		nullUserId sql.NullInt64
	)
	err := rows.Scan(&comment.Id, &comment.UUID, &comment.PostId, &comment.Author, &comment.Email, &comment.Website, &comment.CreatedAt, &comment.Content, &comment.Approved, &comment.UserAgent, &nullParent, &nullUserId)
	comment.Avatar = utils.Gravatar(comment.Email, "50")
	comment.Parent = nullParent.Int64
	comment.UserId = nullUserId.Int64
	return err
}

func GetCommentById(id int64) (*Comment, error) {
	comment := new(Comment)
	row := db.QueryRow(stmtGetCommentById, id)
	err := scanComment(row, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func GetCommentByPostId(id int64) ([]*Comment, error) {
	rows, err := db.Query(stmtGetApprovedCommentListByPostId, id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return extractComments(rows)
}

func DeleteComment(id int64) error {
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtDeleteCommentById, id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}
