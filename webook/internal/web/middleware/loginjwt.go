package middleware

import (
	"github.com/aerbibabo/MyWayToGo/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
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
	return func(ctx *gin.Context) {
		for _, path := range l.ignorePaths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		authStr := ctx.GetHeader("Authorization")
		if authStr == "" {
			println("error 1")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		authSegments := strings.Split(authStr, "Bearer ")
		tokenStr := authSegments[1]
		if tokenStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := &web.UserClime{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(web.JWTKey), nil
		})
		if err != nil {
			println("error 2" + err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil {
			println("error 2")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("claims", claims)
	}
}
