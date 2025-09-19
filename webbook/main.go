package main

import (
	"github.com/aerbibabo/MyWayToGo/webbook/internal/repostory"
	"github.com/aerbibabo/MyWayToGo/webbook/internal/repostory/dao"
	"github.com/aerbibabo/MyWayToGo/webbook/internal/service"
	"github.com/aerbibabo/MyWayToGo/webbook/internal/web"
	"github.com/aerbibabo/MyWayToGo/webbook/internal/web/middleware"
	"github.com/gin-contrib/sessions"
	//"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

func main() {
	db := initDB()

	err := dao.InitTables(db)
	if err != nil {
		panic(err.Error())
	}

	server := initWebServer()

	//登录校验
	//store := cookie.NewStore([]byte("secret"))
	redisStore, err := redis.NewStore(16, "tcp", "localhost:6379",
		"root", "",
		[]byte("0apssj7hhulymyb0"),
		[]byte("hdfpa75c6x1y00k9"))
	if err != nil {
		panic("redis连接失败" + err.Error())
		return
	}

	server.Use(sessions.Sessions("ssid", redisStore))

	server.Use(middleware.NewLoginMiddleware().
		IgnorePath("/users/login").
		IgnorePath("/users/signup").Build())

	u := initUser(db)
	u.RegisterRoutes(server)

	err = server.Run(":8080")
	if err != nil {
		panic(err.Error())
		return
	}
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	ur := repostory.NewUserRepository(ud)
	us := service.NewUserService(ur)
	c := web.NewUserHandler(us)
	return c
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			println("[cors]请求来源" + origin)
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "https://company.name")
		},
		MaxAge: 86400,
	}))
	return server
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/webook"))
	if err != nil {
		panic(err.Error())
	}

	return db
}
