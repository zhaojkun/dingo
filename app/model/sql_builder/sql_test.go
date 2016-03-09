package SQL

import (
	"testing"
)

func TestSQLBuilder(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{
			Select(`*`).From(`posts`).Where(`status="published"`).OrderBy("posts.published_at DESC").Limit(`5`).Offset(`0`).SQL(),
			`SELECT * FROM posts WHERE status="published" ORDER BY posts.published_at DESC LIMIT 5 OFFSET 0`,
		},
		{
			Select(`*`).From(`posts`).SQL(),
			`SELECT * FROM posts`,
		},
		{
			Select(`*`).From(`posts, posts_tags`).Where(`posts_tags.post_id = posts.id`, `posts_tags.tag_id = ?`, `page = 0`, `status = 'published'`).OrderBy(`posts.published_at DESC`).Limit(`?`).Offset(`?`).SQL(),
			`SELECT * FROM posts, posts_tags WHERE posts_tags.post_id = posts.id AND posts_tags.tag_id = ? AND page = 0 AND status = 'published' ORDER BY posts.published_at DESC LIMIT ? OFFSET ?`,
		},
		{
			Select(`*`).From(`posts`).OrderBy("published_at").SQL(),
			`SELECT * FROM posts ORDER BY published_at`,
		},
		{
			Select(`*`).From(`posts`).SQL(),
			`SELECT * FROM posts`,
		},
	}

	for _, c := range cases {
		if c.in != c.out {
			t.Errorf("SQL builder error: \nExpected %v \nGot      %v", c.out, c.in)
		}
	}
}
