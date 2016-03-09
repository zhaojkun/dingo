package model

import (
	"fmt"
	"github.com/dinever/dingo/app/utils"
	"github.com/dinever/golf"
)

var tokens map[string]*Token

type Token struct {
	Value      string
	UserId     int64
	CreateTime int64
	ExpireTime int64
}

func CreateToken(u *User, context *Golf.Context, expire int64) *Token {
	t := new(Token)
	t.UserId = u.Id
	t.CreateTime = utils.NowUnix()
	t.ExpireTime = t.CreateTime + expire
	t.Value = utils.Sha1(fmt.Sprintf("%s-%s-%d-%d", context.Request.RemoteAddr, context.Request.UserAgent(), t.CreateTime, t.UserId))
	tokens[t.Value] = t
	return t
}

// get token by token value.
func GetTokenByValue(v string) *Token {
	return tokens[v]
}

// get tokens of given user.
func GetTokensByUser(u *User) []*Token {
	ts := make([]*Token, 0)
	for _, t := range tokens {
		if t.UserId == u.Id {
			ts = append(ts, t)
		}
	}
	return ts
}

// remove a token by token value.
func RemoveToken(v string) {
	delete(tokens, v)
}

// clean all expired tokens in memory.
// do not write to json.
func CleanTokens() {
	for k, t := range tokens {
		if !t.IsValid() {
			delete(tokens, k)
		}
	}
}

func (t *Token) IsValid() bool {
	user, _ := GetUserById(t.UserId)
	if user == nil {
		return false
	}
	return t.ExpireTime > utils.NowUnix()
}
