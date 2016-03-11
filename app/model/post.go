package model

import (
	"database/sql"
	"github.com/dinever/dingo/app/utils"
	"github.com/twinj/uuid"
	"log"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	Id              int64
	UUID            string
	Title           string
	Slug            string
	Markdown        string
	Html            string
	Image           string
	CommentNum      int64
	Comments        []*Comment
	IsFeatured      bool
	IsPublished     bool
	status          string
	IsPage          bool
	AllowComment    bool
	Category        string
	Hits            int64
	Language        string
	MetaTitle       string
	MetaDescription string
	Author          *User
	userId          int64
	CreatedAt       *time.Time
	CreatedBy       int64
	UpdatedAt       *time.Time
	UpdatedBy       int64
	PublishedAt     *time.Time
	PublishedBy     int64
	Tags            []*Tag
}

func NewPost() *Post {
	post := new(Post)
	post.UUID = uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)
	var createdAt time.Time
	createdAt = time.Now()
	post.CreatedAt = &createdAt
	return post
}

func (p *Post) TagString() string {
	var tagString string
	for i, t := range p.Tags {
		if i != len(p.Tags)-1 {
			tagString += t.Name + ", "
		} else {
			tagString += t.Name
		}
	}
	return tagString
}

func (p *Post) Url() string {
	return "/" + p.Slug
}

func (p *Post) Summary() string {
	text := strings.Split(p.Markdown, "<!--more-->")[0]
	return utils.Markdown2Html(text)
}

func (p *Post) Excerpt() string {
	return utils.Html2Excerpt(p.Html, 255)
}

func (p *Post) Save() error {
	if p.Id == 0 {
		// Insert post
		postId, err := InsertPost(p)
		if err != nil {
			return err
		}
		p.Id = postId
	} else {
		err := UpdatePost(p)
		if err != nil {
			return err
		}
	}
	tagIds := make([]int64, 0)
	// Insert tags
	for _, t := range p.Tags {
		var createdAt time.Time
		createdAt = time.Now()
		t.CreatedAt = &createdAt
		t.CreatedBy = p.userId
		t.Hidden = !p.IsPublished
		t.Save()
		tagIds = append(tagIds, t.Id)
	}
	// Delete old post-tag projections
	err := DeletePostTagsByPostId(p.Id)
	// Insert postTags
	if err != nil {
		return err
	}
	for _, tagId := range tagIds {
		err := InsertPostTag(p.Id, tagId)
		if err != nil {
			return err
		}
	}
	return DeleteOldTags()
}

func InsertPost(p *Post) (int64, error) {
	if !PostChangeSlug(p.Slug) {
		p.Slug = generateNewSlug(p.Slug, 1)
	}
	if p.IsPublished {
		p.status = "published"
	} else {
		p.status = "draft"
	}
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	var result sql.Result
	if p.IsPublished {
		result, err = writeDB.Exec(stmtInsertPost, nil, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), p.Title, p.Slug, p.Markdown, p.Html, p.IsFeatured, p.IsPage, p.AllowComment, p.status, p.Image, p.CreatedBy, p.CreatedAt, p.CreatedBy, p.CreatedAt, p.CreatedBy, p.CreatedAt, p.CreatedBy)
	} else {
		result, err = writeDB.Exec(stmtInsertPost, nil, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), p.Title, p.Slug, p.Markdown, p.Html, p.IsFeatured, p.IsPage, p.AllowComment, p.status, p.Image, p.CreatedBy, p.CreatedAt, p.CreatedBy, p.CreatedAt, p.CreatedBy, nil, nil)
	}
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	postId, err := result.LastInsertId()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	return postId, writeDB.Commit()
}

func InsertPostTag(post_id int64, tag_id int64) error {
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertPostTag, nil, post_id, tag_id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func UpdatePost(p *Post) error {
	currentPost, err := GetPostById(p.Id)
	if err != nil {
		return err
	}
	if p.Slug != currentPost.Slug && !PostChangeSlug(p.Slug) {
		p.Slug = generateNewSlug(p.Slug, 1)
	}
	status := "draft"
	if p.IsPublished {
		status = "published"
	}
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// If the updated post is published for the first time, add publication date and user
	if p.IsPublished && !currentPost.IsPublished {
		_, err = writeDB.Exec(stmtUpdatePostPublished, p.Title, p.Slug, p.Markdown, p.Html, p.IsFeatured, p.IsPage, p.AllowComment, status, p.Image, p.CreatedAt, p.CreatedBy, p.CreatedAt, p.CreatedBy, p.Id)
	} else {
		_, err = writeDB.Exec(stmtUpdatePost, p.Title, p.Slug, p.Markdown, p.Html, p.IsFeatured, p.IsPage, p.AllowComment, status, p.Image, p.CreatedAt, p.CreatedBy, p.Id)
	}
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func DeletePostTagsByPostId(post_id int64) error {
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtDeletePostTagsByPostId, post_id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func DeletePostById(id int64) error {
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtDeletePostById, id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	err = writeDB.Commit()
	if err != nil {
		return err
	}
	err = DeletePostTagsByPostId(id)
	if err != nil {
		return err
	}
	return DeleteOldTags()
}

func GetPostCreationDateById(post_id int64) (*time.Time, error) {
	var date time.Time
	// Get number of posts
	row := db.QueryRow(stmtGetPostCreationDateById, post_id)
	err := row.Scan(&date)
	if err != nil {
		return &date, err
	}
	return &date, nil
}

func GetPostById(id int64) (*Post, error) {
	// Get post
	row := db.QueryRow(stmtGetPostById, id)
	return extractPost(row)
}

func GetPostBySlug(slug string) (*Post, error) {
	// Get post
	row := db.QueryRow(stmtGetPostBySlug, slug)
	return extractPost(row)
}

func GetAllPostsByTag(tagId int64) ([]*Post, error) {
	// Get posts
	rows, err := db.Query(stmtGetAllPostsByTag, tagId)
	defer rows.Close()
	if err != nil {
		log.Printf("[Error] Can not get posts from tag: %v", err.Error())
		return nil, err
	}
	posts, err := extractPosts(rows)
	if err != nil {
		log.Printf("[Error] Can not scan posts from tag: %v", err.Error())
		return nil, err
	}
	return posts, nil
}

func GetNumberOfPosts(isPage bool, published bool) (int64, error) {
	var count int64
	selector := postCountSelector.Copy()
	if published {
		selector.Where(`status = "published"`)
	}
	if isPage {
		selector.Where(`page = 1`)
	} else {
		selector.Where(`page = 0`)
	}
	var row *sql.Row
	row = db.QueryRow(selector.SQL())
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetPostList(page, size int64, isPage bool, onlyPublished bool, orderBy string) ([]*Post, *utils.Pager, error) {
	var (
		pager *utils.Pager
	)
	count, err := GetNumberOfPosts(isPage, onlyPublished)
	pager = utils.NewPager(page, size, count)
	selector := postSelector.Copy()
	if onlyPublished {
		selector.Where(`status = "published"`)
	}
	if isPage {
		selector.Where(`page = 1`)
	} else {
		selector.Where(`page = 0`)
	}
	selector.OrderBy(orderBy)
	// Get posts
	rows, err := db.Query(selector.Limit(`?`).Offset(`?`).SQL(), size, pager.Begin-1)
	defer rows.Close()
	if err != nil {
		log.Printf("[Error]: ", err.Error())
		return nil, nil, err
	}
	posts, err := extractPosts(rows)
	if err != nil {
		return nil, nil, err
	}
	return posts, pager, nil
}

func GetAllPostList(isPage bool, onlyPublished bool, orderBy string) ([]*Post, error) {
	selector := postSelector.Copy()
	if onlyPublished {
		selector.Where(`status = "published"`)
	}
	if isPage {
		selector.Where(`page = 1`)
	} else {
		selector.Where(`page = 0`)
	}
	selector.OrderBy(orderBy)
	// Get posts
	rows, err := db.Query(selector.SQL())
	defer rows.Close()
	if err != nil {
		log.Printf("[Error]: ", err.Error())
		return nil, err
	}
	posts, err := extractPosts(rows)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func scanPost(rows Row, post *Post) error {
	// TODO: CommentNum
	post.CommentNum = 0
	var (
		nullImage       sql.NullString
		nullUpdatedBy   sql.NullInt64
		nullPublishedBy sql.NullInt64
	)
	err := rows.Scan(&post.Id, &post.UUID, &post.Title, &post.Slug, &post.Markdown,
		&post.Html, &post.IsFeatured, &post.IsPage, &post.AllowComment, &post.CommentNum, &post.status, &nullImage,
		&post.userId, &post.CreatedAt, &post.CreatedBy, &post.UpdatedAt, &nullUpdatedBy, &post.PublishedAt, &nullPublishedBy)
	post.Image = nullImage.String
	return err
}

func extractPosts(rows *sql.Rows) ([]*Post, error) {
	posts := make([]*Post, 0)
	for rows.Next() {
		post := new(Post)
		err := scanPost(rows, post)
		if err != nil {
			return nil, err
		}
		// If there was no publication date attached to the post, make its creation date the date of the post
		if post.PublishedAt == nil {
			post.PublishedAt, err = GetPostCreationDateById(post.Id)
			if err != nil {
				return nil, err
			}
		}
		// Evaluate status
		if post.status == "published" {
			post.IsPublished = true
		} else {
			post.IsPublished = false
		}
		// Get user
		post.Author, err = GetUserById(post.userId)
		if err != nil {
			return nil, err
		}
		// Get tags
		post.Tags, err = GetTags(post.Id)
		if err != nil {
			return nil, err
		}
		post.Comments, err = GetCommentByPostId(post.Id)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func extractPost(row *sql.Row) (*Post, error) {
	post := new(Post)
	err := scanPost(row, post)
	if err != nil {
		return nil, err
	}
	// If there was no publication date attached to the post, make its creation date the date of the post
	if post.PublishedAt == nil {
		post.PublishedAt, err = GetPostCreationDateById(post.Id)
		if err != nil {
			return nil, err
		}
	}
	// Evaluate status
	if post.status == "published" {
		post.IsPublished = true
	} else {
		post.IsPublished = false
	}
	// Get user
	post.Author, err = GetUserById(post.userId)
	if err != nil {
		return nil, err
	}
	// Get tags
	post.Tags, err = GetTags(post.Id)
	if err != nil {
		return nil, err
	}
	// Get comments
	post.Comments, err = GetCommentByPostId(post.Id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func PostChangeSlug(slug string) bool {
	_, err := GetPostBySlug(slug)
	if err != nil {
		return true
	}
	return false
}

func generateNewSlug(slug string, suffix int) string {
	newSlug := slug + "-" + strconv.Itoa(suffix)
	if !PostChangeSlug(newSlug) {
		return generateNewSlug(slug, suffix+1)
	}
	return newSlug
}
