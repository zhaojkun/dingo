package handler

import (
	"github.com/dinever/dingo/app/model"
)

func getAllPosts() []*model.Post {
	posts, _ := model.GetAllPostList(false, true, "published_at")
	return posts
}

func getAllTags() []*model.Tag {
	tags, _ := model.GetAllTags()
	return tags
}

func getRecentPosts() []*model.Post {
	posts, _, _ := model.GetPostList(1, 5, false, true, "published_at")
	return posts
}

