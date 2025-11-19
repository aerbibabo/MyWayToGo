package web

import (
	"errors"
	"github.com/aerbibabo/MyWayToGo/webook/internal/domain"
	"github.com/aerbibabo/MyWayToGo/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"time"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	JWTKey               = "0apssj7hhulymyb0"
)

type UserHandler struct {
	svc              *service.UserService
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
}

// 对外暴露实例化函数
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:              svc,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (user *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile", user.Profile)
	ug.POST("/login", user.Login)
	ug.POST("/edit", user.Edit)
	ug.POST("/signup", user.SignUp)
}

func (user *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	//接收json参数并使用结构体解析 绑定变量
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "error[10000]")
		return
	}
	//校验邮箱
	matchString, err := user.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "error[10001]")
		return
	}
	if !matchString {
		ctx.String(http.StatusOK, "邮箱格式不正确")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不一样")
		return
	}

	//校验密码格式
	isPassword, err := user.passwordRegexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "error[10002]")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码格式不正确 数字 特殊字符 8位以上")
		return
	}

	err = user.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突, 请换一个")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "服务器异常,注册失败"+err.Error())
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (user *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "error[10000]")
		return
	}

	if req.Email == "" || req.Password == "" {
		ctx.String(http.StatusOK, "email或密码为空")
	}

	loginUser, err := user.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//sess := sessions.Default(ctx)
	//sess.Set("userId", loginUser.Id)
	//sess.Options(sessions.Options{
	//	Path:   "/",
	//	MaxAge: 30,
	//})
	//_ = sess.Save()

	claims := UserClime{
		Userid: loginUser.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	//使用jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(JWTKey))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "生成jwt错误")
		return
	}

	ctx.Header("x-jwt-token", tokenStr)

	ctx.String(http.StatusOK, "登陆成功")

}

func (user *UserHandler) Edit(ctx *gin.Context) {
	ctx.String(http.StatusOK, "edit")

}
func (user *UserHandler) Profile(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	if !ok {
		println(ok)

		ctx.String(http.StatusInternalServerError, "profile error")
		return
	}

	claims, ok := c.(*UserClime)
	if !ok {
		println(ok)
		ctx.String(http.StatusInternalServerError, "profile error2")
	}

	ctx.String(http.StatusOK, "profile"+strconv.FormatInt(claims.Userid, 10))
}

type UserClime struct {
	jwt.RegisteredClaims
	Userid int64
}
