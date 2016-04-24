package handler

import (
	"github.com/dinever/dingo/app/model"
	"github.com/dinever/dingo/app/utils"
	"github.com/dinever/golf"
	"github.com/twinj/uuid"
	"strconv"
	"time"
)

func AdminHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	ctx.Loader("admin").Render("home.html", map[string]interface{}{
		"Title":    "Dashboard",
		"Statis":   model.NewStatis(ctx.App),
		"User":     user,
		"Messages": model.GetUnreadMessages(),
		"Monitor":  utils.ReadMemStats(),
	})
}

func ProfileHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	ctx.Loader("admin").Render("profile.html", map[string]interface{}{
		"Title": "Profile",
		"User":  user,
	})
}

func ProfileChangeHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	if user.Email != ctx.Request.FormValue("email") && !model.UserChangeEmail(ctx.Request.FormValue("email")) {
		ctx.JSON(map[string]interface{}{"res": false, "msg": "A user with that email address already exists."})
		return
	}
	user.Name = ctx.Request.FormValue("name")
	user.Slug = ctx.Request.FormValue("slug")
	user.Email = ctx.Request.FormValue("email")
	user.Avatar = utils.Gravatar(ctx.Request.FormValue("email"), "180")
	user.Website = ctx.Request.FormValue("url")
	user.Bio = ctx.Request.FormValue("bio")
	err := user.UpdateUser(user.Id)
	if err != nil {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": err.Error(),
		})
	}
	ctx.JSON(map[string]interface{}{"res": true})
}

func PostCreateHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	c := model.NewPost()
	ctx.Loader("admin").Render("edit_article.html", map[string]interface{}{
		"Title": "New Post",
		"Post":  c,
		"User":  user,
	})
}

func PostSaveHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	p := model.NewPost()
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)
	p.Id = int64(idInt)
	p.Title = ctx.Request.FormValue("title")
	p.Slug = ctx.Request.FormValue("slug")
	p.Markdown = ctx.Request.FormValue("content")
	p.Html = utils.Markdown2Html(p.Markdown)
	p.Tags = model.GenerateTagsFromCommaString(ctx.Request.FormValue("tag"))
	p.AllowComment = ctx.Request.FormValue("comment") == "on"
	p.Category = ctx.Request.FormValue("category")
	p.CreatedBy = user.Id
	p.UpdatedBy = user.Id
	p.IsPublished = ctx.Request.FormValue("status") == "on"
	p.IsPage = false
	p.Author = user
	p.Hits = 1
	var e error
	e = p.Save()
	if e != nil {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": e.Error()})
		return
	}
	ctx.JSON(map[string]interface{}{
		"res":     true,
		"content": p,
	})
}

func AdminPostHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	i, _ := strconv.Atoi(ctx.Request.FormValue("page"))
	articles, pager, err := model.GetPostList(int64(i), 10, false, false, "created_at DESC")
	if err != nil {
		panic(err)
	}
	ctx.Loader("admin").Render("posts.html", map[string]interface{}{
		"Title": "Posts",
		"Posts": articles,
		"User":  user,
		"Pager": pager,
	})
}

func PostEditHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	id := ctx.Param("id")
	articleId, _ := strconv.Atoi(id)
	c, err := model.GetPostById(int64(articleId))
	if c == nil || err != nil {
		ctx.Redirect("/admin/posts/")
		return
	}
	ctx.Loader("admin").Render("edit_article.html", map[string]interface{}{
		"Title": "Edit Post",
		"Post":  c,
		"User":  user,
	})
}

func PostRemoveHandler(ctx *golf.Context) {
	id := ctx.Param("id")
	articleId, _ := strconv.Atoi(id)
	err := model.DeletePostById(int64(articleId))
	if err != nil {
		ctx.JSON(map[string]interface{}{
			"res": false,
		})
	} else {
		ctx.JSON(map[string]interface{}{
			"res": true,
		})
	}
}

func PageCreateHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	c := model.NewPost()
	ctx.Loader("admin").Render("edit_article.html", map[string]interface{}{
		"Title": "New Page",
		"Post":  c,
		"User":  user,
	})
}

func AdminPageHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	i, _ := strconv.Atoi(ctx.Request.FormValue("page"))
	pages, pager, err := model.GetPostList(int64(i), 10, true, false, `created_at`)
	println(pages)
	if err != nil {
		panic(err)
	}
	ctx.Loader("admin").Render("pages.html", map[string]interface{}{
		"Title": "Pages",
		"Pages": pages,
		"User":  user,
		"Pager": pager,
	})
}

func PageSaveHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	p := model.NewPost()
	p.Id = 0
	if !model.PostChangeSlug(ctx.Request.FormValue("slug")) {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": "The slug of this post has conflicts with another post."})
		return
	}
	p.Title = ctx.Request.FormValue("title")
	p.Slug = ctx.Request.FormValue("slug")
	p.Markdown = ctx.Request.FormValue("content")
	p.Html = utils.Markdown2Html(p.Markdown)
	p.Tags = model.GenerateTagsFromCommaString(ctx.Request.FormValue("tag"))
	p.AllowComment = ctx.Request.FormValue("comment") == "on"
	p.Category = ctx.Request.FormValue("category")
	p.CreatedBy = user.Id
	p.IsPublished = ctx.Request.FormValue("status") == "on"
	p.IsPage = true
	p.Author = user
	p.Hits = 1
	var e error
	e = p.Save()
	if e != nil {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": e.Error(),
		})
		return
	}
	ctx.JSON(map[string]interface{}{
		"res":     true,
		"content": p,
	})
}

func CommentViewHandler(ctx *golf.Context) {
	i, _ := strconv.Atoi(ctx.Request.FormValue("page"))
	user, _ := ctx.Session.Get("user")
	comments, pager, err := model.GetCommentList(int64(i), 10)
	if err != nil {
		panic(err)
	}
	ctx.Loader("admin").Render("comments.html", map[string]interface{}{
		"Title":    "Comments",
		"Comments": comments,
		"User":     user,
		"Pager":    pager,
	})
}

func CommentAddHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	pid, _ := strconv.Atoi(ctx.Request.FormValue("pid"))
	parent, err := model.GetCommentById(int64(pid))
	if err != nil {
		panic(err)
	}
	comment := new(model.Comment)
	comment.Author = user.Name
	comment.Email = user.Email
	comment.Website = user.Website
	comment.Content = ctx.Request.FormValue("content")
	comment.Avatar = utils.Gravatar(comment.Email, "50")
	comment.Parent = parent.Id
	comment.PostId = parent.PostId
	comment.Ip = ctx.Request.RemoteAddr
	comment.UserAgent = ctx.Request.UserAgent()
	comment.UserId = user.Id
	comment.Approved = true
	t := time.Now()
	comment.CreatedAt = &t
	id, err := comment.Save()
	if err != nil {
		panic(err)
	}
	comment.Id = id
	ctx.JSON(map[string]interface{}{
		"res":     true,
		"comment": comment.ToJson(),
	})
	model.CreateMessage("comment", comment)
}

func CommentUpdateHandler(ctx *golf.Context) {
	id, _ := strconv.Atoi(ctx.Request.FormValue("id"))
	c, err := model.GetCommentById(int64(id))
	if err != nil {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": err.Error(),
		})
	}
	c.Approved = true
	c.Save()
	ctx.JSON(map[string]interface{}{
		"res": true,
	})
}

func CommentRemoveHandler(ctx *golf.Context) {
	id, _ := strconv.Atoi(ctx.Request.FormValue("id"))
	err := model.DeleteComment(int64(id))
	if err != nil {
		ctx.JSON(map[string]interface{}{
			"res": true,
			"msg": err.Error(),
		})
	}
	ctx.JSON(map[string]interface{}{
		"res": true,
	})
}

func SettingViewHandler(ctx *golf.Context) {
	user, _ := ctx.Session.Get("user")
	ctx.Loader("admin").Render("setting.html", map[string]interface{}{
		"Title":      "Settings",
		"User":       user,
		"Custom":     model.GetCustomSettings(),
		"Navigators": model.GetNavigators(),
	})
}

func SettingUpdateHandler(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	var err error
	for key, value := range ctx.Request.Form {
		setting := new(model.Setting)
		setting.UUID = uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)
		setting.Key = key
		setting.Value = value[0]
		setting.Type = ""
		setting.CreatedBy = user.Id
		now := time.Now()
		setting.CreatedAt = &now
		err = model.SaveSetting(setting)
		if err != nil {
			panic(err)
			ctx.JSON(map[string]interface{}{
				"res": false,
				"msg": err.Error(),
			})
		}
	}
	ctx.JSON(map[string]interface{}{
		"res": true,
	})
}

func SettingCustomHandler(ctx *golf.Context) {
	keys := ctx.Request.Form["key"]
	values := ctx.Request.Form["value"]
	for i, k := range keys {
		if len(k) < 1 {
			continue
		}
		model.SetSetting(k, values[i], "custom")
	}
	ctx.JSON(map[string]interface{}{
		"res": true,
	})
}

func SettingNavHandler(ctx *golf.Context) {
	labels := ctx.Request.Form["label"]
	urls := ctx.Request.Form["url"]
	model.SetNavigators(labels, urls)
	ctx.JSON(map[string]interface{}{
		"res": true,
	})
}

func AdminPasswordPage(ctx *golf.Context) {
	user, _ := ctx.Session.Get("user")
	ctx.Loader("admin").Render("password.html", map[string]interface{}{
		"Title": "Change Password",
		"User":  user,
	})
}

func AdminPasswordChange(ctx *golf.Context) {
	userObj, _ := ctx.Session.Get("user")
	user := userObj.(*model.User)
	oldPassword := ctx.Request.FormValue("old")
	if !user.CheckPassword(oldPassword) {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": "Old password incorrect.",
		})
		return
	}
	newPassword := ctx.Request.FormValue("new")
	user.ChangePassword(newPassword)
	ctx.JSON(map[string]interface{}{
		"res": true,
	})
}

func AdminMonitorPage(ctx *golf.Context) {
	user, _ := ctx.Session.Get("user")
	ctx.Loader("admin").Render("monitor.html", map[string]interface{}{
		"Title":   "Monitor",
		"User":    user,
		"Monitor": utils.ReadMemStats(),
	})
}
