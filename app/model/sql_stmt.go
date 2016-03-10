package model

import (
	"github.com/dinever/dingo/app/model/sql_builder"
)

const schema = `CREATE TABLE IF NOT EXISTS
posts (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  uuid				varchar(36) NOT NULL,
  title varchar(150) NOT NULL,
  slug varchar(150) NOT NULL,
  markdown text,
  html text,
  image text,
  featured tinyint NOT NULL DEFAULT '0',
  page tinyint NOT NULL DEFAULT '0',
  allow_comment tinyint NOT NULL DEFAULT '0',
  comment_num integer NOT NULL DEFAULT '0',
  status varchar(150) NOT NULL DEFAULT 'draft',
  language varchar(6) NOT NULL DEFAULT 'en_US',
  meta_title varchar(150),
  meta_description varchar(200),
  author_id integer NOT NULL,
  created_at datetime NOT NULL,
  created_by integer NOT NULL,
  updated_at datetime,
  updated_by integer,
  published_at datetime,
  published_by integer
);

CREATE TABLE IF NOT EXISTS
	users (
		id					integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid				varchar(36) NOT NULL,
		name				varchar(150) NOT NULL,
		slug				varchar(150) NOT NULL,
		password			varchar(60) NOT NULL,
		email				varchar(254) NOT NULL,
		image				text,
		cover				text,
		bio					varchar(200),
		website				text,
		location			text,
		accessibility		text,
		status				varchar(150) NOT NULL DEFAULT 'active',
		language			varchar(6) NOT NULL DEFAULT 'en_US',
		meta_title			varchar(150),
		meta_description	varchar(200),
		last_login			datetime,
		created_at			datetime NOT NULL,
		created_by			integer NOT NULL,
		updated_at			datetime,
		updated_by			integer
	);

CREATE TABLE IF NOT EXISTS
	categories (
		id					integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid				varchar(36) NOT NULL,
		name				varchar(150) NOT NULL,
		slug				varchar(150) NOT NULL,
		description			varchar(200),
		parent_id			integer,
		meta_title			varchar(150),
		meta_description	varchar(200),
		created_at			datetime NOT NULL,
		created_by			integer NOT NULL,
		updated_at			datetime,
		updated_by			integer
	);

CREATE TABLE IF NOT EXISTS
	tags (
		id					integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid				varchar(36) NOT NULL,
		name				varchar(150) NOT NULL,
		slug				varchar(150) NOT NULL,
		description			varchar(200),
    image       text,
    hidden      boolean NOT NULL DEFAULT 0,
		parent_id			integer,
		meta_title			varchar(150),
		meta_description	varchar(200),
		created_at			datetime NOT NULL,
		created_by			integer NOT NULL,
		updated_at			datetime,
		updated_by			integer
	);

	CREATE TABLE IF NOT EXISTS
	comments (
		id					integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid				varchar(36) NOT NULL,
		post_id				varchar(150) NOT NULL,
		author				varchar(150) NOT NULL,
		author_email				varchar(150) NOT NULL,
		author_url varchar(200) NOT NULL,
		author_ip varchar(100) NOT NULL,
		created_at datetime NOT NULL,
		content text NOT NULL,
		approved tinyint NOT NULL DEFAULT '0',
		agent varchar(255) NOT NULL,
		type varchar(20),
		parent integer,
		user_id integer
	);

CREATE TABLE IF NOT EXISTS
	posts_tags (
		id		integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		post_id	integer NOT NULL,
		tag_id	integer NOT NULL
	);

CREATE TABLE IF NOT EXISTS
	posts_categories (
		id		integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		post_id	integer NOT NULL,
		category_id	integer NOT NULL
	);

CREATE TABLE IF NOT EXISTS
	settings (
		id			integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid		varchar(36) NOT NULL,
		key			varchar(150) NOT NULL,
		value		text,
		type		varchar(150) NOT NULL DEFAULT 'core',
		created_at	datetime NOT NULL,
		created_by	integer NOT NULL,
		updated_at	datetime,
		updated_by	integer
	);

CREATE TABLE IF NOT EXISTS
	roles (
		id			integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid		varchar(36) NOT NULL,
		name		varchar(150) NOT NULL,
		description	varchar(200),
		created_at	datetime NOT NULL,
		created_by	integer NOT NULL,
		updated_at	datetime,
		updated_by	integer
	);
`

// Posts
var postCountSelector = SQL.Select(`count(*)`).From(`posts`)
var stmtGetPublishedPostsCount = postCountSelector.Copy().Where(`status = "published"`).SQL()
var stmtGetAllPostsCount = postCountSelector.Copy().SQL()
var stmtGetPostsCountByUser = postCountSelector.Copy().Where(`author_id = ?`).SQL()
var stmtGetPostsCountByTag = postCountSelector.Copy().Where(`posts_tags.post_id = posts.id`, `posts_tags.tag_id = ?`, `status = 'published'`).SQL()

var postSelector = SQL.Select(`id, uuid, title, slug, markdown, html, featured, page, allow_comment, comment_num, status, image, author_id, created_at, created_by, updated_at, updated_by, published_at, published_by`).From(`posts`)
var stmtGetPublishedPostList = postSelector.Copy().Where(`status = "published"`).OrderBy(`published_at DESC`).Limit(`?`).Offset(`?`).SQL()
var stmtGetAllPostList = postSelector.Copy().OrderBy(`published_at DESC`).Limit(`?`).Offset(`?`).SQL()
var stmtGetPostsByUser = postSelector.Copy().Where(`status = 'published'`, `author_id = ?`).OrderBy(`published_at DESC`).Limit(`?`).Offset(`?`).SQL()

var stmtGetPostById = postSelector.Copy().Where(`id = ?`).SQL()
var stmtGetPostBySlug = postSelector.Copy().Where(`slug = ?`).SQL()

var postsTagsSelector = SQL.Select(`posts.id, posts.uuid, posts.title, posts.slug, posts.markdown, posts.html, posts.featured, posts.page, post.allow_comment, post.comment_num, posts.status, posts.image, posts.author_id, posts.created_at, posts.created_by, posts.updated_at, posts.updated_by, posts.published_at, posts.published_by`).From(`posts, posts_tags`)
var stmtGetPostsByTag = postsTagsSelector.Copy().Where(`status = 'published'`, `posts_tags.post_id = posts_id`, `posts_tags.tag_id = ?`).OrderBy(`published_at DESC`).Limit(`?`).Offset(`?`).SQL()
var stmtGetAllPostsByTag = postsTagsSelector.Copy().Where(`posts_tags.post_id = posts.id`, `posts_tags.tag_id = ?`).OrderBy(`published_at DESC`).SQL()

var pageCountSelector = SQL.Select(`count(*)`).From(`posts`).Where(`page = 1`)
var stmtGetPublishedPagesCount = pageCountSelector.Copy().Where(`status = "published"`).SQL()
var stmtGetAllPagesCount = pageCountSelector.Copy().SQL()
var stmtGetPagesCountByUser = pageCountSelector.Copy().Where(`author_id = ?`).SQL()
var stmtGetPagesCountByTag = pageCountSelector.Copy().Where(`posts_tags.post_id = posts.id`, `posts_tags.tag_id = ?`, `status = 'published'`).SQL()

var pageSelector = SQL.Select(`id, uuid, title, slug, markdown, html, featured, page, status, image, author_id, created_at, created_by, updated_at, updated_by, published_at, published_by`).From(`posts`).Where(`page = 1`)

// Comments

var commentCountSelector = SQL.Select(`count(*)`).From(`comments`)
var stmtGetAllCommentCount = commentCountSelector.SQL()
var commentSelector = SQL.Select(`id, uuid, post_id, author, author_email, author_url, created_at, content, approved, agent, parent, user_id`).From(`comments`)
var stmtGetAllCommentList = commentSelector.Copy().OrderBy(`created_at DESC`).Limit(`?`).Offset(`?`).SQL()
var stmtGetApprovedCommentList = commentSelector.Copy().Where(`approved = 1`).OrderBy(`created_at DESC`).Limit(`?`).Offset(`?`).SQL()
var stmtGetCommentById = commentSelector.Copy().Where(`id = ?`).SQL()
var stmtInsertComment = `INSERT OR REPLACE INTO comments (id, uuid, post_id, author, author_email, author_url, author_ip, created_at, content, approved, agent, parent, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
var stmtDeleteCommentById = `DELETE FROM comments WHERE id = ?`

// Users
const stmtGetUserById = `SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE id = ?`
const stmtGetUserBySlug = `SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE slug = ?`
const stmtGetUserByName = `SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE name = ?`
const stmtGetUserByEmail = `SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE email = ?`
const stmtGetHashedPasswordByEmail = `SELECT password FROM users WHERE email = ?`
const stmtGetUsersCount = `SELECT count(*) FROM users`
const stmtGetUsersCountByEmail = `SELECT count(*) FROM users where email = ?`

// Tags
const stmtGetAllTags = `SELECT id, name, slug FROM tags`
const stmtGetTags = `SELECT tag_id FROM posts_tags WHERE post_id = ?`
const stmtGetTagById = `SELECT id, name, slug FROM tags WHERE id = ?`
const stmtGetTagBySlug = `SELECT id, name, slug, hidden FROM tags WHERE slug = ?`
const stmtGetTagIdBySlug = `SELECT id FROM tags WHERE slug = ?`

// Settings
const stmtGetBlog = `SELECT value FROM settings WHERE key = ?`
const stmtGetPostCreationDateById = `SELECT created_at FROM posts WHERE id = ?`

const stmtInsertPost = `INSERT INTO posts (id, uuid, title, slug, markdown, html, featured, page, allow_comment, status, image, author_id, created_at, created_by, updated_at, updated_by, published_at, published_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
const stmtInsertUser = `INSERT INTO users (id, uuid, name, slug, password, email, image, cover, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
const stmtInsertRoleUser = `INSERT INTO roles_users (id, role_id, user_id) VALUES (?, ?, ?)`
const stmtInsertTag = `INSERT INTO tags (id, uuid, name, slug, created_at, created_by, updated_at, updated_by, hidden) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
const stmtInsertPostTag = `INSERT INTO posts_tags (id, post_id, tag_id) VALUES (?, ?, ?)`
const stmtInsertSetting = `INSERT INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

const stmtUpdatePost = `UPDATE posts SET title = ?, slug = ?, markdown = ?, html = ?, featured = ?, page = ?, allow_comment = ?, status = ?, image = ?, updated_at = ?, updated_by = ? WHERE id = ?`
const stmtUpdatePostPublished = `UPDATE posts SET title = ?, slug = ?, markdown = ?, html = ?, featured = ?, page = ?, status = ?, image = ?, updated_at = ?, updated_by = ?, published_at = ?, published_by = ? WHERE id = ?`
const stmtUpdateSettings = `UPDATE settings SET value = ?, updated_at = ?, updated_by = ? WHERE key = ?`
const stmtUpdateUser = `UPDATE users SET name = ?, slug = ?, email = ?, image = ?, cover = ?, bio = ?, website = ?, location = ?, updated_at = ?, updated_by = ? WHERE id = ?`
const stmtUpdateLastLogin = `UPDATE users SET last_login = ? WHERE id = ?`
const stmtUpdateUserPassword = `UPDATE users SET password = ?, updated_at = ?, updated_by = ? WHERE id = ?`
const stmtUpdateTag = `UPDATE tags SET uuid = ?, name = ?, slug =?, updated_at = ?, updated_by = ?, hidden = ? WHERE id = ?`

const stmtDeletePostTagsByPostId = `DELETE FROM posts_tags WHERE post_id = ?`
const stmtDeletePostById = `DELETE FROM posts WHERE id = ?`
const stmtDeleteOldTags = `DELETE FROM tags WHERE id IN (SELECT id FROM tags EXCEPT SELECT tag_id FROM posts_tags)`
