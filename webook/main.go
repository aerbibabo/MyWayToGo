package main

import (
	"github.com/aerbibabo/MyWayToGo/webook/config"
	"github.com/aerbibabo/MyWayToGo/webook/internal/repostory"
	"github.com/aerbibabo/MyWayToGo/webook/internal/repostory/dao"
	"github.com/aerbibabo/MyWayToGo/webook/internal/service"
	"github.com/aerbibabo/MyWayToGo/webook/internal/web"
	"github.com/aerbibabo/MyWayToGo/webook/internal/web/middleware"
	"github.com/aerbibabo/MyWayToGo/webook/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/sessions"
	sessredis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/contrib/cors"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

func main() {
	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "Hello World!")
	//})
	//err := server.Run(":8080")
	//if err != nil {
	//	panic(err.Error())
	//	return
	//}

	db := initDB()

	//初始化数据库 , 笨方法, 以后不用
	//err1 := dao.InitTables(db)
	//if err1 != nil {
	//	panic(err1.Error())
	//}

	server := initWebServer()

	u := initUser(db)
	u.RegisterRoutes(server)

	err := server.Run(":8080")
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
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"X-Jwt-Token"},
		AllowOriginFunc: func(origin string) bool {
			println("[cors]请求来源" + origin)
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "https://company.name")
		},
		MaxAge: 86400,
	}))

	//登录校验
	redisStore, err := sessredis.NewStore(16, "tcp", config.Config.Redis.Addr,
		"root",
		"",
		[]byte("0apssj7hhulymyb0"),
		[]byte("hdfpa75c6x1y00k9"))
	if err != nil {
		panic("redis连接失败" + err.Error())
	}

	//注册session store
	server.Use(sessions.Sessions("ssid", redisStore))

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
		DB:   config.Config.Redis.DB,
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(middleware.NewLoginMiddleware().
		IgnorePath("/users/login").
		IgnorePath("/users/signup").Build())
	return server
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	return db
}
