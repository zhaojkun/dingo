package handler

import (
	"github.com/dinever/golf"
	"github.com/dinever/dingo/app/model"
	"strconv"
)

func AuthMiddleware(next golf.HandlerFunc) golf.HandlerFunc {
	fn := func(ctx *golf.Context) {
		//		user, _ := model.GetUserByEmail("dingpeixuan911@gmail.com")
		//		ctx.Data["user"] = user
		//		next(ctx)
		userNum, err := model.GetNumberOfUsers()
		if err == nil {
			if userNum == 0 {
				ctx.Redirect("/signup/")
				return
			}
		}
		tokenStr, err := ctx.Request.Cookie("token-value")
		if err == nil {
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
				ctx.Session.Set("user", user)
				next(ctx)
			} else {
				ctx.Redirect("/login/")
			}
		} else {
			ctx.Redirect("/login/")
		}
	}
	return fn
}
