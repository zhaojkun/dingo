package Dingo

import (
	"fmt"
	"github.com/dinever/dingo/app/handler"
	"github.com/dinever/dingo/app/model"
	"github.com/dinever/dingo/app/utils"
	"github.com/dinever/golf"
	"os"
	"path/filepath"
)

var (
	App *golf.Application
)

func Init() {
	App = golf.New()
	err := model.Initialize("dingo.db")
	if err != nil {
		panic(err)
	}

	App.Config.Set("app/static_dir", "static")
	App.Config.Set("app.log_dir", "tmp/log")
	App.Config.Set("app/upload_dir", "upload")
	upload_dir, _ := App.Config.GetString("app/upload_dir", "upload")
	registerMiddlewares()
	registerFuncMap()
	handler.RegisterFunctions(App)
	theme := model.GetSettingValue("theme")
	App.View.SetTemplateLoader("base", "view")
	App.View.SetTemplateLoader("admin", filepath.Join("view", "admin"))
	App.View.SetTemplateLoader("theme", filepath.Join("view", theme))
	//	static_dir, _ := App.Config.GetString("app/static_dir", "static")
	App.Static("/upload/", upload_dir)
	App.Static("/", filepath.Join("view", "admin", "assets"))
	App.Static("/", filepath.Join("view", theme, "assets"))

	App.SessionManager = golf.NewMemorySessionManager()
	App.Error(404, handler.NotFoundHandler)
}

func registerFuncMap() {
	App.View.FuncMap["DateFormat"] = utils.DateFormat
	App.View.FuncMap["DateInt64"] = utils.DateInt64
	App.View.FuncMap["DateString"] = utils.DateString
	App.View.FuncMap["DateTime"] = utils.DateTime
	App.View.FuncMap["Now"] = utils.Now
	App.View.FuncMap["Html2Str"] = utils.Html2Str
	App.View.FuncMap["FileSize"] = utils.FileSize
	App.View.FuncMap["Setting"] = model.GetSettingValue
	App.View.FuncMap["Navigator"] = model.GetNavigators
	App.View.FuncMap["Md2html"] = utils.Markdown2HtmlTemplate
}

func registerMiddlewares() {
	App.Use(
		golf.LoggingMiddleware(os.Stdout),
		golf.RecoverMiddleware,
		golf.SessionMiddleware,
	)
}

func CreateSampleData() {
	model.SetSetting("site_url", "http://example.com/", "blog")
	model.SetSetting("title", "Dingo Blog", "blog")
	model.SetSetting("sub_title", "Another blog created by Dingo", "blog")
}

func registerAdminURLHandlers() {
	authChain := golf.NewChain(handler.AuthMiddleware)
	App.Get("/login/", handler.AuthLoginPageHandler)
	App.Post("/login/", handler.AuthLoginHandler)

	App.Get("/signup/", handler.AuthSignUpPageHandler)
	App.Post("/signup/", handler.AuthSignUpHandler)

	App.Get("/logout/", handler.AuthLogoutHandler)

	App.Get("/admin/", authChain.Final(handler.AdminHandler))

	App.Get("/admin/profile/", authChain.Final(handler.ProfileHandler))
	App.Post("/admin/profile/", authChain.Final(handler.ProfileChangeHandler))

	App.Get("/admin/editor/post/", authChain.Final(handler.PostCreateHandler))
	App.Post("/admin/editor/post/", authChain.Final(handler.PostSaveHandler))

	App.Get("/admin/editor/page/", authChain.Final(handler.PageCreateHandler))
	App.Post("/admin/editor/page/", authChain.Final(handler.PageSaveHandler))

	App.Get("/admin/posts/", authChain.Final(handler.AdminPostHandler))
	App.Get("/admin/editor/:id/", authChain.Final(handler.PostEditHandler))
	App.Post("/admin/editor/:id/", authChain.Final(handler.PostSaveHandler))
	App.Delete("/admin/editor/:id/", authChain.Final(handler.PostRemoveHandler))

	App.Get("/admin/pages/", authChain.Final(handler.AdminPageHandler))

	App.Get("/admin/comments/", authChain.Final(handler.CommentViewHandler))
	App.Post("/admin/comments/", authChain.Final(handler.CommentAddHandler))
	App.Put("/admin/comments/", authChain.Final(handler.CommentUpdateHandler))
	App.Delete("/admin/comments/", authChain.Final(handler.CommentRemoveHandler))

	App.Get("/admin/setting/", authChain.Final(handler.SettingViewHandler))
	App.Post("/admin/setting/", authChain.Final(handler.SettingUpdateHandler))
	App.Post("/admin/setting/custom/", authChain.Final(handler.SettingCustomHandler))
	App.Post("/admin/setting/nav/", authChain.Final(handler.SettingNavHandler))
	//
	App.Get("/admin/files/", authChain.Final(handler.FileViewHandler))
	App.Delete("/admin/files/", authChain.Final(handler.FileRemoveHandler))
	App.Post("/admin/files/upload/", authChain.Final(handler.FileUploadHandler))

	App.Get("/admin/password/", authChain.Final(handler.AdminPasswordPage))
	App.Post("/admin/password/", authChain.Final(handler.AdminPasswordChange))

	App.Get("/admin/monitor/", authChain.Final(handler.AdminMonitorPage))
}

func registerHomeHandler() {
	statsChain := golf.NewChain()
	App.Get("/", statsChain.Final(handler.HomeHandler))
	App.Get("/page/:page/", handler.HomeHandler)
	App.Post("/comment/:id/", handler.CommentHandler)
	App.Get("/tag/:tag/", handler.TagHandler)
	App.Get("/tag/:tag/page/:page/", handler.TagHandler)
	App.Get("/feed/", handler.RssHandler)
	App.Get("/sitemap.xml", handler.SiteMapHandler)
	App.Get("/:slug/", statsChain.Final(handler.ContentHandler))
}

func Run(portNumber string) {
	registerAdminURLHandlers()
	registerHomeHandler()
	fmt.Println("Application Started on port " + portNumber)
	App.Run(":" + portNumber)
}
