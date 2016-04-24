package handler

import (
	"github.com/dinever/dingo/app/model"
	"github.com/dinever/golf"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func FileViewHandler(ctx *golf.Context) {
	user, _ := ctx.Session.Get("user")
	uploadDir, _ := ctx.App.Config.GetString("upload_dir", "upload")
	uploadDir = path.Clean(uploadDir)
	dir, err := ctx.Query("dir")
	dir = path.Clean(dir)
	var (
		ParentDir  string
		IsChildDir bool
	)
	if err == nil && dir != uploadDir {
		IsChildDir = true
		ParentDir = path.Join(dir, "..")
	} else {
		IsChildDir = false
		dir = uploadDir
	}
	var files []*model.File
	if model.CheckSafe(dir, uploadDir) {
		files = model.GetFileList(dir)
	} else {
		ctx.Abort(403)
		return
	}
	ctx.Loader("admin").Render("files.html", map[string]interface{}{
		"Title":      "Files",
		"Files":      files,
		"User":       user,
		"CurrentDir": dir,
		"IsChildDir": IsChildDir,
		"ParentDir":  ParentDir,
	})
}

func FileRemoveHandler(ctx *golf.Context) {
	p := ctx.Request.FormValue("path")
	uploadDir, _ := ctx.App.Config.GetString("upload_dir", "upload")
	if model.CheckSafe(p, uploadDir) {
		err := model.RemoveFile(p)
		if err != nil {
			panic(err)
		}
	} else {
		ctx.Abort(403)
		return
	}
	ctx.JSON(map[string]interface{}{
		"res": true,
	})
}

func FileUploadHandler(ctx *golf.Context) {
	req := ctx.Request
	req.ParseMultipartForm(32 << 20)
	f, h, e := req.FormFile("file")
	if e != nil {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": e.Error(),
		})
		return
	}
	data, _ := ioutil.ReadAll(f)
	maxSize, _ := ctx.App.Config.GetInt("app.upload_size", 1024*1024*10)
	defer func() {
		f.Close()
		data = nil
		h = nil
	}()
	if len(data) >= maxSize {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": "File size should be smaller than 10MB.",
		})
		return
	}
	fileExt, _ := ctx.App.Config.GetString("app.upload_files", ".jpg,.png,.gif,.zip,.txt,.doc,.docx,.xls,.xlsx,.ppt,.pptx")
	if !strings.Contains(fileExt, path.Ext(h.Filename)) {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": "Only supports documents, images and zip files.",
		})
		return
	}
	uploadDir, _ := ctx.App.Config.GetString("upload_dir", "upload")
	Url := model.CreateFilePath(uploadDir, h.Filename)
	e = ioutil.WriteFile(Url, data, os.ModePerm)
	if e != nil {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": e.Error(),
		})
		return
	}
	ctx.JSON(map[string]interface{}{
		"res": true,
		"file": map[string]interface{}{
			"url":  Url,
			"name": h.Filename,
		},
	})
}
