package handler

import (
	"github.com/dinever/dingo/app/model"
	"github.com/dinever/golf"
	"regexp"
	"strconv"
	"fmt"
)

const Email string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

var rxEmail = regexp.MustCompile(Email)

func AuthLoginPageHandler(ctx *Golf.Context) {
	ctx.Loader("admin").Render("login.html")
}

func AuthSignUpPageHandler(ctx *Golf.Context) {
	userNum, err := model.GetNumberOfUsers()
	if err != nil {
		ctx.Abort(404)
		return
	}
	if userNum == 0 {
		ctx.Loader("admin").Render("signup.html", nil)
	} else {
		ctx.Abort(404)
		return
	}
}

func AuthSignUpHandler(ctx *Golf.Context) {
	userNum, err := model.GetNumberOfUsers()
	if err != nil || userNum != 0 {
		ctx.Abort(403)
		return
	}

	email := ctx.Request.FormValue("email")
	if !rxEmail.MatchString(email) {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": "Invalid email address.",
		})
		return
	}
	name := ctx.Request.FormValue("name")
	password := ctx.Request.FormValue("password")
	if len(password) < 5 {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": "Password too short.",
		})
		return
	}
	if len(password) > 20 {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": "Password too long.",
		})
		return
	}
	rePassword := ctx.Request.FormValue("re-password")
	if password != rePassword {
		ctx.JSON(map[string]interface{}{
			"res": false,
			"msg": "Password does not match.",
		})
		return
	}
	err = model.CreateNewUser(email, name, password)
	if err != nil {
		ctx.Abort(500)
		return
	}
	user, err := model.GetUserByEmail(email)
	if err != nil {
		ctx.Abort(500)
		return
	}
	rememberMe := ctx.Request.FormValue("remember-me")
	var (
		exp int
		s   *model.Token
	)
	if rememberMe == "on" {
		exp = 3600 * 24 * 3
		s = model.CreateToken(user, ctx, int64(exp))
	} else {
		exp = 0
		s = model.CreateToken(user, ctx, 3600)
	}
	ctx.SetCookie("token-user", strconv.Itoa(int(s.UserId)), exp)
	ctx.SetCookie("token-value", s.Value, exp)
	ctx.JSON(map[string]interface{}{
		"res": true,
	})
}

func AuthLoginHandler(ctx *Golf.Context) {
	email := ctx.Request.FormValue("email")
	password := ctx.Request.FormValue("password")
	rememberMe := ctx.Request.FormValue("remember-me")
	user, err := model.GetUserByEmail(email)
	if user == nil || err != nil {
		ctx.JSON(map[string]interface{}{"res": false})
		return
	}
	if !user.CheckPassword(password) {
		ctx.JSON(map[string]interface{}{"res": false})
		return
	}
	var (
		exp int
		s   *model.Token
	)
	if rememberMe == "on" {
		exp = 3600 * 24 * 3
		s = model.CreateToken(user, ctx, int64(exp))
	} else {
		exp = 0
		s = model.CreateToken(user, ctx, 3600)
	}
	ctx.SetCookie("token-user", strconv.Itoa(int(s.UserId)), exp)
	ctx.SetCookie("token-value", s.Value, exp)
	ctx.JSON(map[string]interface{}{"res": true})
}

func AuthLogoutHandler(ctx *Golf.Context) {
	ctx.SetCookie("token-user", "", -3600)
	ctx.SetCookie("token-value", "", -3600)
	ctx.Redirect("/login/")
}

func verifyUser(ctx *Golf.Context) bool {
	tokenStr, err := ctx.Request.Cookie("token-value")
	if err == nil {
		token := model.GetTokenByValue(tokenStr.Value)
		if token != nil && token.IsValid() {
			return true
		}
	}
	return false
}

func AuthMiddleware(next Golf.Handler) Golf.Handler {
	fn := func(ctx *Golf.Context) {
		//		user, _ := model.GetUserByEmail("dingpeixuan911@gmail.com")
		//		ctx.Data["user"] = user
		//		next(ctx)
		tokenStr, err := ctx.Request.Cookie("token-value")
		if err == nil {
			userNum, err := model.GetNumberOfUsers()
			if err == nil {
				if userNum == 0 {
					ctx.Redirect("/signup/")
					return
				}
			}
			token := model.GetTokenByValue(tokenStr.Value)
			if token != nil && token.IsValid() {
				tokenUser, err := ctx.Request.Cookie("token-user")
				if err != nil {
					panic(err)
				}
				uid, _ := strconv.Atoi(tokenUser.Value)
				user, err := model.GetUserById(int64(uid))
				if err != nil {
					panic(err)
				}
				ctx.Data["user"] = user
				next(ctx)
			} else {
				ctx.App.NotFoundHandler(ctx)
			}
		} else {
			ctx.App.NotFoundHandler(ctx)
		}
	}
	return fn
}
