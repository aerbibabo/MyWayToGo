package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	ignorePaths []string
}

func NewLoginMiddleware() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	l.ignorePaths = append(l.ignorePaths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Time{})
	return func(ctx *gin.Context) {
		for _, path := range l.ignorePaths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//if ctx.Request.URL.Path == "/users/login" ||
		//	ctx.Request.URL.Path == "/users/signup" {
		//	return
		//}

		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		sess.Options(sessions.Options{
			Path:   "/",
			MaxAge: 30,
		})
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		updateTime := sess.Get("update_time")
		if updateTime == nil {
			sess.Set("update_time", now)
			err := sess.Save()
			if err != nil {
				panic("错误原因" + err.Error())
			}
			return
		}

		updateTimeVal, ok := updateTime.(time.Time)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		println(updateTimeVal.Unix())
		if now.Sub(updateTimeVal) > 10*time.Second {
			sess.Set("update_time", now)
			err := sess.Save()
			if err != nil {
				panic("错误原因2" + err.Error())
			}
		}

	}
}
